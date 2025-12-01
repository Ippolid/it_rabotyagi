import os
import csv
import random
from typing import Optional
from fastapi import FastAPI, File, UploadFile, Form, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from groq import Groq
import tempfile
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(title="Interview Trainer API", version="1.0.0")

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure appropriately in production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Groq client will be initialized lazily
client = None

def get_groq_client():
    """Get or initialize Groq client"""
    global client
    if client is None:
        groq_api_key = os.getenv("GROQ_API_KEY")
        if not groq_api_key:
            raise HTTPException(status_code=500, detail="GROQ_API_KEY not configured")
        client = Groq(api_key=groq_api_key)
    return client

# Load questions from CSV
questions_db = []

def load_questions():
    """Load questions from CSV file"""
    global questions_db
    csv_path = os.path.join(os.path.dirname(__file__), "voprosy.last.csv")
    
    if not os.path.exists(csv_path):
        logger.warning(f"Questions file not found: {csv_path}")
        return
    
    try:
        with open(csv_path, "r", encoding="utf-8") as f:
            reader = csv.DictReader(f, delimiter=";")
            questions_db = []
            for idx, row in enumerate(reader, start=1):
                questions_db.append({
                    "id": idx,
                    "question": row.get("Question", "").strip(),
                    "ideal_answer": row.get("Answer", "").strip()
                })
        logger.info(f"Loaded {len(questions_db)} questions from CSV")
    except Exception as e:
        logger.error(f"Error loading questions: {e}")

# Load questions on startup
load_questions()

# Models
class TranscriptionResponse(BaseModel):
    text: str
    language: Optional[str] = None
    duration: Optional[float] = None

class EvaluationRequest(BaseModel):
    question: str
    ideal_answer: str
    user_answer: str
    language: str = "ru"

class EvaluationResponse(BaseModel):
    score: int  # 0-100
    feedback: str
    strengths: list[str]
    improvements: list[str]
    overall_comment: str

class QuestionResponse(BaseModel):
    id: int
    question: str
    ideal_answer: str

@app.get("/")
async def root():
    """Health check endpoint"""
    return {
        "service": "Interview Trainer API",
        "status": "running",
        "version": "1.0.0",
        "questions_loaded": len(questions_db)
    }

@app.get("/api/v1/questions/random", response_model=QuestionResponse)
async def get_random_question():
    """
    Get a random question from the questions database
    
    Returns:
        Random question with its ID and ideal answer
    """
    if not questions_db:
        raise HTTPException(status_code=404, detail="No questions available")
    
    question = random.choice(questions_db)
    logger.info(f"Selected random question ID: {question['id']}")
    
    return QuestionResponse(**question)

@app.get("/api/v1/questions/{question_id}", response_model=QuestionResponse)
async def get_question_by_id(question_id: int):
    """
    Get a specific question by ID
    
    Args:
        question_id: ID of the question
    
    Returns:
        Question with its ID and ideal answer
    """
    if not questions_db:
        raise HTTPException(status_code=404, detail="No questions available")
    
    question = next((q for q in questions_db if q["id"] == question_id), None)
    if not question:
        raise HTTPException(status_code=404, detail=f"Question {question_id} not found")
    
    logger.info(f"Retrieved question ID: {question_id}")
    return QuestionResponse(**question)

@app.get("/api/v1/questions", response_model=list[QuestionResponse])
async def get_all_questions():
    """
    Get all available questions
    
    Returns:
        List of all questions
    """
    if not questions_db:
        raise HTTPException(status_code=404, detail="No questions available")
    
    return [QuestionResponse(**q) for q in questions_db]

@app.post("/api/v1/transcribe", response_model=TranscriptionResponse)
async def transcribe_audio(
    audio: UploadFile = File(...),
    language: Optional[str] = Form(None)
):
    """
    Transcribe audio file using Groq Whisper API
    
    Args:
        audio: Audio file (supported formats: mp3, mp4, mpeg, mpga, m4a, wav, webm, ogg)
        language: Optional language code (e.g., 'ru', 'en')
    
    Returns:
        Transcription text and metadata
    """
    try:
        # Validate file type (check both content_type and extension)
        allowed_types = ["audio/mpeg", "audio/mp4", "audio/wav", "audio/webm", "audio/ogg", "audio/x-m4a", "application/octet-stream"]
        allowed_extensions = [".mp3", ".mp4", ".mpeg", ".mpga", ".m4a", ".wav", ".webm", ".ogg"]
        
        file_extension = os.path.splitext(audio.filename)[1].lower() if audio.filename else ""
        
        # Accept if either content_type is allowed OR extension is allowed
        if audio.content_type not in allowed_types and file_extension not in allowed_extensions:
            raise HTTPException(
                status_code=400,
                detail=f"Unsupported audio format: {audio.content_type} (extension: {file_extension}). Allowed extensions: {', '.join(allowed_extensions)}"
            )
        
        # Save uploaded file to temporary location with proper extension
        suffix = file_extension if file_extension in allowed_extensions else ".webm"
        with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as temp_file:
            content = await audio.read()
            temp_file.write(content)
            temp_file_path = temp_file.name
        
        try:
            # Transcribe using Groq Whisper
            groq_client = get_groq_client()
            with open(temp_file_path, "rb") as audio_file:
                transcription = groq_client.audio.transcriptions.create(
                    file=(audio.filename, audio_file.read()),
                    model="whisper-large-v3-turbo",
                    temperature=0,
                    response_format="verbose_json",
                    language=language
                )
            
            logger.info(f"Transcription completed: {transcription.text[:50]}...")
            
            return TranscriptionResponse(
                text=transcription.text,
                language=transcription.language if hasattr(transcription, 'language') else None,
                duration=transcription.duration if hasattr(transcription, 'duration') else None
            )
        
        finally:
            # Clean up temporary file
            if os.path.exists(temp_file_path):
                os.unlink(temp_file_path)
    
    except Exception as e:
        logger.error(f"Transcription error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Transcription failed: {str(e)}")

@app.post("/api/v1/evaluate", response_model=EvaluationResponse)
async def evaluate_answer(request: EvaluationRequest):
    """
    Evaluate user's interview answer using LLM
    
    Args:
        question: Interview question
        ideal_answer: Expected/ideal answer
        user_answer: User's transcribed answer
        language: Language for feedback (default: 'ru')
    
    Returns:
        Evaluation with score, feedback, strengths, and improvement suggestions
    """
    try:
        # Construct evaluation prompt
        prompt = f"""Ты - опытный технический интервьюер. Оцени ответ кандидата на вопрос собеседования.

ВОПРОС:
{request.question}

ИДЕАЛЬНЫЙ ОТВЕТ:
{request.ideal_answer}

ОТВЕТ КАНДИДАТА:
{request.user_answer}

Проанализируй ответ и предоставь:
1. Оценку от 0 до 100 (где 100 - идеальный ответ)
2. Сильные стороны ответа (2-3 пункта)
3. Области для улучшения (2-3 пункта)
4. Общий комментарий

Формат ответа (строго следуй этому формату):
Оценка: [число от 0 до 100]
Что сказал правильно:
- [сильная сторона 1]
- [сильная сторона 2]
- [сильная сторона 3]
Что улучшить:
- [улучшение 1]
- [улучшение 2]
- [улучшение 3]
Рекомендация:
[общий комментарий на 3-5 предложений]"""

        # Call LLM for evaluation
        groq_client = get_groq_client()
        chat_completion = groq_client.chat.completions.create(
            messages=[
                {
                    "role": "system",
                    "content": "Ты опытный технический интервьюер, который дает конструктивную и честную оценку ответов кандидатов."
                },
                {
                    "role": "user",
                    "content": prompt
                }
            ],
            model="llama-3.3-70b-versatile",  # Fast and good quality
            temperature=0.3,
            max_tokens=2000
        )
        
        response_text = chat_completion.choices[0].message.content
        logger.info(f"Evaluation completed for question: {request.question[:50]}...")
        logger.info(f"LLM Response:\n{response_text}")
        
        # Parse response
        import re
        score = 0
        strengths = []
        improvements = []
        overall_comment = ""
        
        lines = response_text.split('\n')
        current_section = None
        
        for line in lines:
            line = line.strip()
            
            # Parse score
            if line.startswith("Оценка:") or line.startswith("SCORE:"):
                try:
                    score_text = line.split(":")[1].strip()
                    # Extract first number from string
                    numbers = re.findall(r'\d+', score_text)
                    if numbers:
                        score = int(numbers[0])
                        # Ensure score is 0-100
                        score = max(0, min(score, 100))
                except Exception as e:
                    logger.warning(f"Failed to parse score: {e}")
                    score = 50
                continue
            
            # Detect sections - не пропускаем строку после определения секции
            if "правильно" in line.lower() or "strengths" in line.lower():
                current_section = "strengths"
                continue
            elif "улучшить" in line.lower() or "improvements" in line.lower():
                current_section = "improvements"
                continue
            elif "рекомендация" in line.lower() or "comment" in line.lower():
                current_section = "comment"
                continue
            
            # Skip empty lines
            if not line:
                continue
            
            # Parse list items
            if line.startswith("- ") or line.startswith("• "):
                item = line[2:].strip()
                if current_section == "strengths":
                    strengths.append(item)
                elif current_section == "improvements":
                    improvements.append(item)
            # Parse comment - собираем все строки в секции комментария
            elif current_section == "comment":
                overall_comment += line + " "
        
        overall_comment = overall_comment.strip()
        
        # Log parsed results
        logger.info(f"Parsed - Score: {score}, Strengths: {len(strengths)}, Improvements: {len(improvements)}, Comment length: {len(overall_comment)}")
        
        # Fallback: use full response if parsing failed
        if not overall_comment:
            logger.warning("Failed to parse overall_comment, using full response")
            overall_comment = response_text
        
        if not strengths:
            logger.warning("No strengths parsed")
            strengths = ["См. полный анализ в рекомендации"]
        
        if not improvements:
            logger.warning("No improvements parsed")
            improvements = ["См. полный анализ в рекомендации"]
        
        return EvaluationResponse(
            score=score,
            feedback=response_text,
            strengths=strengths,
            improvements=improvements,
            overall_comment=overall_comment
        )
    
    except Exception as e:
        logger.error(f"Evaluation error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Evaluation failed: {str(e)}")

@app.post("/api/v1/interview/complete")
async def complete_interview(
    audio: UploadFile = File(...),
    question_id: Optional[int] = Form(None),
    question: Optional[str] = Form(None),
    ideal_answer: Optional[str] = Form(None),
    language: Optional[str] = Form("ru")
):
    """
    Complete interview flow: transcribe + evaluate in one call
    
    Args:
        audio: Audio file with user's answer
        question_id: Optional question ID to fetch from database
        question: Interview question (required if question_id not provided)
        ideal_answer: Expected answer (required if question_id not provided)
        language: Language code
    
    Returns:
        Combined transcription and evaluation results
    """
    try:
        # Get question and answer from database if question_id provided
        if question_id is not None:
            q = next((q for q in questions_db if q["id"] == question_id), None)
            if not q:
                raise HTTPException(status_code=404, detail=f"Question {question_id} not found")
            question = q["question"]
            ideal_answer = q["ideal_answer"]
            logger.info(f"Using question from DB: ID {question_id}")
        elif not question or not ideal_answer:
            raise HTTPException(
                status_code=400,
                detail="Either question_id or both question and ideal_answer must be provided"
            )
        
        # Step 1: Transcribe
        transcription_result = await transcribe_audio(audio, language)
        
        # Step 2: Evaluate
        evaluation_request = EvaluationRequest(
            question=question,
            ideal_answer=ideal_answer,
            user_answer=transcription_result.text,
            language=language
        )
        evaluation_result = await evaluate_answer(evaluation_request)
        
        return {
            "question_id": question_id,
            "question": question,
            "transcription": transcription_result,
            "evaluation": evaluation_result
        }
    
    except Exception as e:
        logger.error(f"Complete interview error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Interview processing failed: {str(e)}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
