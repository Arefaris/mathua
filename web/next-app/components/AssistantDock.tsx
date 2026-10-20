'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { askAssistant, getConfig, type AssistantLink, type AssistantTurn } from '../lib/api'
import { ArrowUp, Close, MessagesSquare } from './icons'

// The product assistant, as a dock in the corner.
//
// It is mounted once in the root layout rather than per page, so the conversation
// survives navigation within the app: a visitor who asks "where do I start" and
// then follows a link can come back to the answer.
//
// Three constraints are deliberate.
//
// **It always renders, and the server decides what it can do.** `ai_enabled` is
// the one fact `/api/config` exposes (ADR-055). Hiding the button when AI is off
// made a switched-off feature indistinguishable from a missing one: the same
// page looked identical whether the assistant was unimplemented, unconfigured or
// broken. With a provider the panel answers; without one it says so once and
// offers the same destinations as links, which is what a visitor wanted anyway.
// It never pretends to answer, and it never leaves an empty box behind.
//
// **Replies are text, never HTML.** The server returns plain prose and a list of
// links whose hrefs it chose itself, so there is no markup to sanitise and no URL
// the model could have written. The links are rendered as Next `Link`s from that
// list, never parsed out of the reply.
//
// **A failed request is not an empty conversation.** A spent quota and an outage
// produce different copy, because "try again later" is wrong advice for a limit
// that resets tomorrow.

type Bubble = AssistantTurn & { links?: AssistantLink[] }

const SUGGESTIONS = ['What is Mathua?', 'Where should I start?']

// AREAS is what the panel offers when the assistant cannot answer. These are
// routes the app owns regardless of AI, so the dock is useful in a deployment
// with no provider — which is every fresh checkout until a key is set.
const AREAS = [
	{ label: 'Start learning', href: '/learn' },
	{ label: 'Take the diagnostic', href: '/diagnose' },
	{ label: 'Plan a goal', href: '/plan' },
	{ label: 'Concept graph', href: '/graph' },
	{ label: 'Glossary', href: '/glossary' },
]

type Status = 'loading' | 'on' | 'off'

function failureCopy(e: unknown): string {
	const message = e instanceof Error ? e.message : ''
	if (message.includes('429')) {
		return 'That is today\u2019s limit for the assistant. The rest of Mathua works as usual.'
	}
	return 'The assistant is unavailable right now. Everything else still works.'
}

export default function AssistantDock() {
	const [status, setStatus] = useState<Status>('loading')
	const [open, setOpen] = useState(false)
	const [turns, setTurns] = useState<Bubble[]>([])
	const [draft, setDraft] = useState('')
	const [busy, setBusy] = useState(false)
	const [error, setError] = useState('')
	const inputRef = useRef<HTMLInputElement>(null)
	const listRef = useRef<HTMLDivElement>(null)

	useEffect(() => {
		let alive = true
		getConfig()
			.then(c => {
				if (alive) setStatus(c?.ai_enabled ? 'on' : 'off')
			})
			.catch(() => {
				// A config read that fails is not a reason to show an assistant that
				// cannot answer, but it is also not a reason to hide a control the
				// page now always offers; the panel says it is unavailable instead.
				if (alive) setStatus('off')
			})
		return () => {
			alive = false
		}
	}, [])

	useEffect(() => {
		if (open) inputRef.current?.focus()
	}, [open])

	useEffect(() => {
		const list = listRef.current
		if (list) list.scrollTop = list.scrollHeight
	}, [turns, busy, error])

	useEffect(() => {
		if (!open) return
		const onKey = (e: KeyboardEvent) => {
			if (e.key === 'Escape') setOpen(false)
		}
		window.addEventListener('keydown', onKey)
		return () => window.removeEventListener('keydown', onKey)
	}, [open])

	const send = useCallback(
		async (text: string) => {
			const message = text.trim()
			if (!message || busy || status !== 'on') return
			setError('')
			setDraft('')
			// Replay what both sides said, so a clarifying question can be answered.
			// The server keeps only these two roles and bounds the count.
			const history: AssistantTurn[] = turns.map(t => ({ role: t.role, content: t.content }))
			setTurns(t => [...t, { role: 'user', content: message }])
			setBusy(true)
			try {
				const res = await askAssistant({ message, history })
				setTurns(t => [...t, { role: 'assistant', content: res.reply, links: res.links }])
			} catch (e) {
				setError(failureCopy(e))
			} finally {
				setBusy(false)
			}
		},
		[busy, status, turns],
	)

	return (
		<div className="fixed right-4 bottom-20 z-40 flex flex-col items-end gap-2 print:hidden lg:bottom-6">
			{open && (
				<section
					role="dialog"
					aria-label="Mathua assistant"
					className="flex h-[min(28rem,70dvh)] w-[min(22rem,calc(100vw-2rem))] flex-col overflow-hidden rounded-lg border border-mathua-border bg-mathua-bg shadow-xl"
				>
					<header className="flex items-center justify-between border-b border-mathua-border px-3 py-2">
						<p className="font-mono text-[11px] tracking-wider text-mathua-muted uppercase">Assistant</p>
						<button
							type="button"
							onClick={() => setOpen(false)}
							className="rounded p-1 text-mathua-muted hover:bg-mathua-surface-elevated hover:text-mathua-primary"
							aria-label="Close assistant"
						>
							<Close size={14} />
						</button>
					</header>

					{status !== 'on' && (
						<div className="flex-1 space-y-3 overflow-y-auto px-3 py-3">
							<p className="text-[13px] leading-relaxed text-mathua-muted">
								{status === 'loading'
									? 'Connecting to the assistant…'
									: 'The assistant is unavailable right now. These go to the same places:'}
							</p>
							{status === 'off' && (
								<div className="flex flex-wrap gap-2">
									{AREAS.map(a => (
										<Link
											key={a.href}
											href={a.href}
											onClick={() => setOpen(false)}
											className="rounded border border-mathua-border px-2 py-1 font-mono text-[11px] text-mathua-blue hover:bg-mathua-surface-elevated"
										>
											{a.label}
										</Link>
									))}
								</div>
							)}
						</div>
					)}

					{status === 'on' && (
						<>
					<div ref={listRef} aria-live="polite" className="flex-1 space-y-3 overflow-y-auto px-3 py-3">
						{turns.length === 0 && (
							<div className="space-y-2">
								<p className="text-[13px] leading-relaxed text-mathua-muted">
									Ask what Mathua is, or where to begin. Answers are short, and you can follow the links.
								</p>
								<div className="flex flex-wrap gap-2">
									{SUGGESTIONS.map(s => (
										<button
											key={s}
											type="button"
											onClick={() => void send(s)}
											className="rounded border border-mathua-border px-2 py-1 font-mono text-[11px] text-mathua-blue hover:bg-mathua-surface-elevated"
										>
											{s}
										</button>
									))}
								</div>
							</div>
						)}

						{turns.map((t, i) => (
							<div key={i} className={t.role === 'user' ? 'flex justify-end' : 'flex justify-start'}>
								<div
									className={[
										'max-w-[85%] rounded-lg px-3 py-2 text-[13px] leading-relaxed whitespace-pre-wrap',
										t.role === 'user'
											? 'bg-mathua-surface-elevated text-mathua-primary'
											: 'border border-mathua-border text-mathua-primary',
									].join(' ')}
								>
									{t.content}
									{t.links && t.links.length > 0 && (
										<div className="mt-2 flex flex-wrap gap-2">
											{t.links.map(l => (
												<Link
													key={l.href}
													href={l.href}
													onClick={() => setOpen(false)}
													className="rounded border border-mathua-border px-2 py-1 font-mono text-[11px] text-mathua-blue hover:bg-mathua-surface-elevated"
												>
													{l.label}
												</Link>
											))}
										</div>
									)}
								</div>
							</div>
						))}

						{busy && <p className="font-mono text-[11px] text-mathua-muted">Thinking…</p>}
						{error && <p className="text-[12px] text-mathua-red">{error}</p>}
					</div>

					<form
						onSubmit={e => {
							e.preventDefault()
							void send(draft)
						}}
						className="flex items-center gap-2 border-t border-mathua-border px-2 py-2"
					>
						<input
							ref={inputRef}
							value={draft}
							onChange={e => setDraft(e.target.value)}
							maxLength={2000}
							placeholder="Ask about Mathua"
							aria-label="Your question"
							className="min-w-0 flex-1 bg-transparent px-1 text-[13px] text-mathua-primary outline-none placeholder:text-mathua-muted"
						/>
						<button
							type="submit"
							disabled={busy || !draft.trim()}
							aria-label="Send"
							className="rounded p-1.5 text-mathua-blue disabled:text-mathua-muted"
						>
							<ArrowUp size={16} />
						</button>
					</form>
						</>
					)}
				</section>
			)}

			<button
				type="button"
				onClick={() => setOpen(v => !v)}
				aria-expanded={open}
				aria-label={open ? 'Close assistant' : 'Open assistant'}
				className="flex h-11 w-11 items-center justify-center rounded-full border border-mathua-border bg-mathua-surface text-mathua-primary shadow-lg hover:bg-mathua-surface-elevated"
			>
				{open ? <Close size={18} /> : <MessagesSquare size={18} />}
			</button>
		</div>
	)
}
