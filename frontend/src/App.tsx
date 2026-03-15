import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { AlertCircle, ArrowRight, Loader2, RefreshCw } from 'lucide-react'
import { useState, useRef, useEffect } from 'react'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8081'

const LANGUAGES = [
  { code: 'en', name: 'English' },
  { code: 'es', name: 'Spanish' },
  { code: 'fr', name: 'French' },
  { code: 'de', name: 'German' },
  { code: 'it', name: 'Italian' },
  { code: 'pt', name: 'Portuguese' },
  { code: 'vi', name: 'Vietnamese' },
  { code: 'ja', name: 'Japanese' },
  { code: 'zh', name: 'Chinese' },
] as const

type Message = {
  id: string
  role: 'user' | 'assistant'
  sourceLang: string
  targetLang: string
  text: string
  translatedText?: string
  error?: string
}

async function translate(
  text: string,
  sourceLang: string,
  targetLang: string
): Promise<{ translated_text: string }> {
  const res = await fetch(`${API_BASE}/api/translate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      text,
      source_lang: sourceLang,
      target_lang: targetLang,
    }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error((data as { error?: string }).error || res.statusText)
  return data as { translated_text: string }
}

function App() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [sourceLang, setSourceLang] = useState('en')
  const [targetLang, setTargetLang] = useState('es')
  const [loading, setLoading] = useState(false)
  const listRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    listRef.current?.scrollTo({ top: listRef.current.scrollHeight, behavior: 'smooth' })
  }, [messages])

  const runTranslation = (trimmed: string, src: string, tgt: string, assistantId: string) => {
    setLoading(true)
    translate(trimmed, src, tgt)
      .then(({ translated_text }) => {
        setMessages((m) =>
          m.map((msg) =>
            msg.id === assistantId ? { ...msg, translatedText: translated_text } : msg
          )
        )
      })
      .catch((err) => {
        setMessages((m) =>
          m.map((msg) =>
            msg.id === assistantId
              ? { ...msg, error: err instanceof Error ? err.message : 'Translation failed' }
              : msg
          )
        )
      })
      .finally(() => setLoading(false))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const trimmed = input.trim()
    if (!trimmed || loading) return

    const userMsg: Message = {
      id: crypto.randomUUID(),
      role: 'user',
      sourceLang,
      targetLang,
      text: trimmed,
    }
    setMessages((m) => [...m, userMsg])
    const assistantId = crypto.randomUUID()
    setMessages((m) => [
      ...m,
      { ...userMsg, id: assistantId, role: 'assistant', text: trimmed },
    ])
    runTranslation(trimmed, sourceLang, targetLang, assistantId)
  }

  const handleRetry = (msg: Message) => {
    if (loading || !msg.error) return
    runTranslation(msg.text, msg.sourceLang, msg.targetLang, msg.id)
    setMessages((m) =>
      m.map((m) => (m.id === msg.id ? { ...m, error: undefined, translatedText: undefined } : m))
    )
  }

  const pairs = (() => {
    const out: { id: string; original: string; assistant: Message }[] = []
    for (let i = 0; i < messages.length - 1; i++) {
      const a = messages[i]
      const b = messages[i + 1]
      if (a.role === 'user' && b.role === 'assistant') {
        out.push({ id: b.id, original: a.text, assistant: b })
        i++
      }
    }
    return out
  })()

  return (
    <div className="app-root bg-[hsl(var(--background))]">
      <main className="app-page">
        <Card className="app-card-fixed w-full flex flex-col flex-1 min-h-0 flex-shrink-0">
          <CardHeader className="app-card-header">
            <h1 className="app-heading">AI Translate</h1>
            <p className="app-caption">
              Type to translate between languages.
            </p>
          </CardHeader>
          <CardContent className="app-card-content app-card-content-cols">
            <aside className="app-input-column">
              <form onSubmit={handleSubmit} className="app-form app-form-vertical">
                <div className="app-lang-row">
                  <Select
                    value={sourceLang}
                    onValueChange={setSourceLang}
                    disabled={loading}
                  >
                    <SelectTrigger className="w-full" aria-busy={loading}>
                      <SelectValue placeholder="From" />
                    </SelectTrigger>
                    <SelectContent>
                      {LANGUAGES.map((l) => (
                        <SelectItem key={l.code} value={l.code}>
                          {l.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <ArrowRight className="app-lang-arrow" aria-hidden />
                  <Select
                    value={targetLang}
                    onValueChange={setTargetLang}
                    disabled={loading}
                  >
                    <SelectTrigger className="w-full" aria-busy={loading}>
                      <SelectValue placeholder="To" />
                    </SelectTrigger>
                    <SelectContent>
                      {LANGUAGES.map((l) => (
                        <SelectItem key={l.code} value={l.code}>
                          {l.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <textarea
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  placeholder="Type or paste text to translate…"
                  disabled={loading}
                  className="app-textarea app-textarea-fill"
                  aria-label="Text to translate"
                  aria-busy={loading}
                  rows={5}
                />
                <Button
                  type="submit"
                  disabled={loading}
                  size="default"
                  className="app-submit-btn"
                  aria-busy={loading}
                >
                  {loading ? (
                    <>
                      <Loader2 className="app-btn-spinner" aria-hidden />
                      <span>Translating…</span>
                    </>
                  ) : (
                    'Translate'
                  )}
                </Button>
              </form>
            </aside>
            <div className="app-column-divider" aria-hidden />
            <section ref={listRef} className="app-output-column space-y-4" aria-label="Translations">
              {pairs.length === 0 && (
                <div className="app-empty-state">
                  <p className="app-empty-title">No translations yet</p>
                  <p className="app-empty-caption">
                    Type text and click Translate. Results appear here.
                  </p>
                </div>
              )}
              {pairs.map(({ id, assistant }) => (
                <article
                  key={id}
                  className={`app-translation-block ${assistant.error ? 'animate-shake' : ''} ${
                    assistant.translatedText ? 'animate-message-in' : ''
                  }`}
                >
                  <div className="app-translation-result">
                    {assistant.error ? (
                      <div className="app-translation-error-block">
                        <div className="app-translation-error-body">
                          <AlertCircle className="app-translation-error-icon" aria-hidden />
                          <div className="app-translation-error-text">
                            <p className="app-translation-error-title">Translation failed</p>
                            <p className="app-translation-error-message">{assistant.error}</p>
                          </div>
                        </div>
                        <Button
                          variant="link"
                          size="sm"
                          className="app-translation-retry"
                          onClick={() => handleRetry(assistant)}
                          disabled={loading}
                        >
                          <RefreshCw className="size-3.5 mr-1" />
                          Retry
                        </Button>
                      </div>
                    ) : assistant.translatedText ? (
                      <p className="app-translation-text">{assistant.translatedText}</p>
                    ) : (
                      <div className="app-loading-state">
                        <Loader2 className="app-loading-spinner" aria-hidden />
                        <span>Translating…</span>
                      </div>
                    )}
                  </div>
                </article>
              ))}
            </section>
          </CardContent>
        </Card>
      </main>
    </div>
  )
}

export default App
