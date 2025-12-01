from __future__ import annotations

from dataclasses import dataclass
from typing import List, Dict, Any, Optional
from pathlib import Path
import re, math, json, os

import numpy as np
import pandas as pd
from tqdm.auto import tqdm

import torch
import torchaudio

# Whisper через transformers
# from transformers import AutoProcessor, AutoModelForSpeechSeq2Seq

# LangChain
from langchain_core.documents import Document
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_core.embeddings import Embeddings
from langchain_community.vectorstores import FAISS

# LLM через llama-cpp + LangChain
from llama_cpp import Llama
from langchain_core.language_models import LLM
from langchain_core.prompts import ChatPromptTemplate

# sentence-transformers
from sentence_transformers import SentenceTransformer

def normalize_ws(s: str) -> str:
    return re.sub(r"\s+", " ", (s or "")).strip()

df = pd.read_csv('/kaggle/input/it-rabotyagi-rag/voprosy.last.csv', on_bad_lines='skip', sep = ';')

class FridaEmbeddings(Embeddings):
    """
    Обёртка над SentenceTransformer (по умолчанию FRIDA).
    Нормализуем эмбеддинги → можно использовать dot-product как косинус.
    """
    def __init__(self, model_name: str = "ai-forever/FRIDA", device: Optional[str] = None):
        self.model = SentenceTransformer(model_name, device=device)

    def embed_documents(self, texts: List[str]) -> List[List[float]]:
        embs = self.model.encode(
            texts,
            convert_to_numpy=True,
            normalize_embeddings=True,
            prompt_name="search_document"
        )
        return embs.tolist()

    def embed_query(self, text: str) -> List[float]:
        emb = self.model.encode(
            [text],
            convert_to_numpy=True,
            normalize_embeddings=True,
            prompt_name="search_query"
        )[0]
        return emb.tolist()

class EmbeddingsWithProgress(Embeddings):
    def __init__(self, base_embeddings: Embeddings, batch_size: int = 256, desc: str = "Embedding"):
        self.base = base_embeddings
        self.batch_size = batch_size
        self.desc = desc

    def embed_query(self, text: str) -> List[float]:
        return self.base.embed_query(text)

    def embed_documents(self, texts: List[str]) -> List[List[float]]:
        out: List[List[float]] = []
        bs = self.batch_size
        total = len(texts)
        iters = math.ceil(total / bs)
        for i in tqdm(range(iters), desc=self.desc, unit="batch"):
            s, e = i * bs, min((i + 1) * bs, total)
            out.extend(self.base.embed_documents(texts[s:e]))
        return out
    
import whisper

whisper_model = whisper.load_model("medium", device="cuda:1")

def transcribe_audio(path: str) -> str:
    result = whisper_model.transcribe(path, fp16=False)
    return result["text"].strip()

from huggingface_hub import login
from langchain_core.language_models import LLM
from typing import Optional, List

login(token="")

llm_cpp = Llama.from_pretrained(
    repo_id="bartowski/gemma-2-9b-it-GGUF",
    filename="gemma-2-9b-it-Q6_K_L.gguf",
    n_ctx=4096,
    n_threads=8,
    n_gpu_layers=-1,
    verbose=False,
)

class LlamaCppLLM(LLM):
    def __init__(self, model: Llama):
        super().__init__()
        self._model = model
    
    def _call(self, prompt: str, stop: Optional[List[str]] = None) -> str:
        output = self._model(
            prompt,
            max_tokens=512,
            temperature=0.3,
            stop=stop or ["</s>", "User:", "\n\n\n"],
        )
        return output["choices"][0]["text"].strip()
    
    @property
    def _llm_type(self) -> str:
        return "llama_cpp"


llm = LlamaCppLLM(model=llm_cpp)

qa_df = df

qa_df = qa_df.rename(columns={
    "Question": "question",
    "Answer": "ideal_answer"
})

assert {"question", "ideal_answer"}.issubset(qa_df.columns)

docs: List[Document] = []
for _, r in qa_df.iterrows():
    q = normalize_ws(str(r["question"]))
    a = normalize_ws(str(r["ideal_answer"]))
    tags = normalize_ws(str(r.get("tags", "")))
    page = f"Вопрос: {q}\nИдеальный ответ: {a}"
    docs.append(Document(
        page_content=page,
        metadata={"question": q, "tags": tags}
    ))

splitter = RecursiveCharacterTextSplitter(
    chunk_size=512,
    chunk_overlap=100,
    separators=["\n\n", "\n", ". ", " ", ""],
)

chunks = splitter.split_documents(docs)
print(f"Документов: {len(docs)}, чанков: {len(chunks)}")

has_cuda = torch.cuda.is_available()
frida_device = "cuda" if has_cuda else None

base_emb = FridaEmbeddings(device=frida_device)
emb_with_progress = EmbeddingsWithProgress(base_emb, batch_size=128, desc="Embed QA")

vector_store = FAISS.from_documents(chunks, embedding=emb_with_progress)
vector_store.save_local("faiss_interview_qa")

# Когда нужно использовать (в том же или в новом ноутбуке):
base_emb = FridaEmbeddings(device=frida_device)
vector_store = FAISS.load_local("faiss_interview_qa", base_emb, allow_dangerous_deserialization=True)

retriever = vector_store.as_retriever(
    search_type="similarity",
    search_kwargs={"k": 5}
)

system_prompt = """
Ты — строгий, но доброжелательный интервьюер по Machine Learning / Backend.
Твоя задача — сравнить ответ кандидата с идеальными ответами и дать развёрнутую обратную связь.

Всегда отвечай строго в формате:

[Оценка: X/5]
[Что кандидат ответил правильно]
[Что кандидат упустил]
[Ошибки в ответе]
[Рекомендации, как улучшить ответ]

Не переписывай идеальный ответ целиком.
Сконцентрируйся на разнице между идеалом и ответом.
"""

human_template = """
Вопрос на собеседовании:
{question}

Ответ кандидата (после транскрибации из аудио):
{user_answer}

Идеальные ответы из базы знаний (несколько вариантов):
{context}

Сгенерируй подробную обратную связь согласно формату.
"""

feedback_prompt = ChatPromptTemplate.from_messages([
    ("system", system_prompt),
    ("human", human_template),
])

feedback_chain = feedback_prompt | llm


def generate_feedback_from_audio(audio_path: str, question_text: str):
    user_answer = transcribe_audio(audio_path)
    print("Распознанный текст ответа:\n", user_answer)
    print("-" * 80)

    # ✔ Правильный вызов retriever
    docs = retriever.invoke(user_answer)

    context = "\n\n---\n\n".join([d.page_content for d in docs])

    raw_feedback = feedback_chain.invoke({
        "question": question_text,
        "user_answer": user_answer,
        "context": context,
    })

    return {
        "transcription": user_answer,
        "feedback_raw": raw_feedback,
    }

audio_file = "/kaggle/input/it-rabotyagi-rag/audio_2025-11-05_22-08-14.ogg"  # путь до ответа
question = "Что такое линейная регрессия?"   # или подтягиваешь из своей системы

result = generate_feedback_from_audio(audio_file, question)

print("\n=== Итоговая обратная связь ===\n")
print(result["feedback_raw"])