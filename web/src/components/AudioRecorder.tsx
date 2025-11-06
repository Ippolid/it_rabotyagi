import { useEffect, useRef, useState } from 'react'

// Simple in-browser recorder placeholder. It DOES NOT call backend yet.
// It records microphone and creates a Blob URL; we "fake" transcript as empty
// and provide hook for future Whisper integration.

export default function AudioRecorder({ onTranscript }: { onTranscript: (text: string) => void }) {
  const mediaRecorderRef = useRef<MediaRecorder | null>(null)
  const chunksRef = useRef<BlobPart[]>([])
  const [recording, setRecording] = useState(false)
  const [audioUrl, setAudioUrl] = useState<string | null>(null)

  useEffect(() => () => {
    if (audioUrl) URL.revokeObjectURL(audioUrl)
  }, [audioUrl])

  async function start() {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      const mr = new MediaRecorder(stream)
      chunksRef.current = []
      mr.ondataavailable = (e) => {
        if (e.data && e.data.size > 0) chunksRef.current.push(e.data)
      }
      mr.onstop = () => {
        const blob = new Blob(chunksRef.current, { type: 'audio/webm' })
        const url = URL.createObjectURL(blob)
        setAudioUrl(url)
        // Placeholder: set empty transcript; actual STT will be added later
        onTranscript('')
      }
      mr.start()
      mediaRecorderRef.current = mr
      setRecording(true)
    } catch (e) {
      console.error('Recorder error', e)
    }
  }

  function stop() {
    mediaRecorderRef.current?.stop()
    setRecording(false)
  }

  return (
    <div className="space-y-3">
      <div className="flex gap-3">
        {!recording ? (
          <button className="btn-primary" onClick={start}>Запись</button>
        ) : (
          <button className="btn-secondary" onClick={stop}>Стоп</button>
        )}
      </div>
      {audioUrl && (
        <audio controls src={audioUrl} className="w-full" />
      )}
    </div>
  )
}
