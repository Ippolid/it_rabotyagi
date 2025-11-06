export type RecordingState = 'idle' | 'recording' | 'stopped'

export class MicRecorder {
  private mediaRecorder?: MediaRecorder
  private chunks: Blob[] = []
  public state: RecordingState = 'idle'

  async start(constraints: MediaStreamConstraints = { audio: true }): Promise<void> {
    if (this.state === 'recording') return
    const stream = await navigator.mediaDevices.getUserMedia(constraints)
    this.chunks = []
    this.mediaRecorder = new MediaRecorder(stream)
    this.mediaRecorder.ondataavailable = (e) => {
      if (e.data && e.data.size > 0) this.chunks.push(e.data)
    }
    this.mediaRecorder.start()
    this.state = 'recording'
  }

  async stop(): Promise<Blob> {
    if (!this.mediaRecorder) throw new Error('Recorder not started')
    if (this.state !== 'recording') throw new Error('Not recording')
    const mr = this.mediaRecorder
    const stopped = new Promise<void>((resolve) => {
      mr.onstop = () => resolve()
    })
    mr.stop()
    this.state = 'stopped'
    await stopped
    const mime = this.chunks[0]?.type || 'audio/webm'
    const blob = new Blob(this.chunks, { type: mime })
    // stop all tracks
    mr.stream.getTracks().forEach(t => t.stop())
    return blob
  }
}

// Optional client-side fallback using Web Speech API (non‑standard, Chrome only)
export async function tryBrowserTranscribe(timeoutMs = 15000): Promise<string | null> {
  const SR: any = (window as any).webkitSpeechRecognition || (window as any).SpeechRecognition
  if (!SR) return null
  return new Promise((resolve) => {
    const r = new SR()
    r.lang = 'ru-RU'
    r.interimResults = false
    let text = ''
    const to = setTimeout(() => { try { r.abort() } catch {} resolve(text || null) }, timeoutMs)
    r.onresult = (e: any) => {
      for (let i = e.resultIndex; i < e.results.length; ++i) {
        if (e.results[i].isFinal) text += e.results[i][0].transcript
      }
    }
    r.onerror = () => { clearTimeout(to); resolve(text || null) }
    r.onend = () => { clearTimeout(to); resolve(text || null) }
    r.start()
  })
}

// Very simple local text similarity (placeholder until server embeddings ready)
export function localTextSimilarity(a: string, b: string): number {
  const tok = (s: string) => s
    .toLowerCase()
    .replace(/[^а-яa-z0-9\s]/gi, ' ')
    .split(/\s+/)
    .filter(Boolean)
  const wa = tok(a)
  const wb = tok(b)
  const set = new Set([...wa, ...wb])
  const vec = (w: string[]) => {
    const m = new Map<string, number>()
    w.forEach(t => m.set(t, (m.get(t) || 0) + 1))
    return m
  }
  const va = vec(wa)
  const vb = vec(wb)
  let dot = 0, na = 0, nb = 0
  set.forEach(t => {
    const x = va.get(t) || 0
    const y = vb.get(t) || 0
    dot += x * y
    na += x * x
    nb += y * y
  })
  if (na === 0 || nb === 0) return 0
  return dot / (Math.sqrt(na) * Math.sqrt(nb))
}
