import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import AssistantDock from '../components/AssistantDock'

const getConfigMock = vi.fn()
const askAssistantMock = vi.fn()

vi.mock('../lib/api', () => ({
	getConfig: (...a: unknown[]) => getConfigMock(...a),
	askAssistant: (...a: unknown[]) => askAssistantMock(...a),
}))

vi.mock('next/link', () => ({
	default: ({ href, children, ...rest }: { href: string; children: React.ReactNode }) => (
		<a href={href} {...rest}>
			{children}
		</a>
	),
}))

// The icon wrapper renders animated SVG through @animateicons/react, which is
// not what this test is about. A plain span keeps the assertions about copy and
// links, so a change to the icon set does not fail a behaviour test.
vi.mock('../components/icons', () => ({
	ArrowUp: () => <span aria-hidden="true">up</span>,
	Close: () => <span aria-hidden="true">close</span>,
	MessagesSquare: () => <span aria-hidden="true">chat</span>,
}))

function openPanel() {
	fireEvent.click(screen.getByRole('button', { name: 'Open assistant' }))
}

beforeEach(() => {
	getConfigMock.mockReset()
	askAssistantMock.mockReset()
})

describe('AssistantDock', () => {
	it('offers the button when the server says AI is off', async () => {
		getConfigMock.mockResolvedValue({ ai_enabled: false })
		render(<AssistantDock />)

		const trigger = await screen.findByRole('button', { name: 'Open assistant' })
		openPanel()

		expect(trigger).toBeInTheDocument()
		expect(await screen.findByText(/assistant is unavailable right now/i)).toBeInTheDocument()
		// The panel still routes a visitor somewhere useful, without a provider.
		expect(screen.getByRole('link', { name: 'Start learning' })).toHaveAttribute('href', '/learn')
		// And it does not offer an input that cannot answer.
		expect(screen.queryByLabelText('Your question')).toBeNull()
	})

	it('still offers the button when the config read fails', async () => {
		getConfigMock.mockRejectedValue(new Error('offline'))
		render(<AssistantDock />)

		expect(await screen.findByRole('button', { name: 'Open assistant' })).toBeInTheDocument()
		openPanel()

		expect(await screen.findByText(/assistant is unavailable right now/i)).toBeInTheDocument()
	})

	it('answers, replays history and renders the reply as text', async () => {
		getConfigMock.mockResolvedValue({ ai_enabled: true })
		askAssistantMock.mockResolvedValue({
			reply: '<b>not markup</b> Learn',
			segments: [
				{ text: '<b>not markup</b> ' },
				{ link: { id: 'learn', label: 'Learn', href: '/learn' } },
			],
		})
		render(<AssistantDock />)

		await screen.findByRole('button', { name: 'Open assistant' })
		openPanel()

		const input = await screen.findByLabelText('Your question')
		fireEvent.change(input, { target: { value: 'Where do I start?' } })
		fireEvent.click(screen.getByRole('button', { name: 'Send' }))

		expect(await screen.findByText('<b>not markup</b>')).toBeInTheDocument()
		expect(screen.queryByText('not markup', { selector: 'b' })).toBeNull()
		// The label is printed once, as the link itself: no chip repeats it under
		// the sentence, which is what made a reply read "start here: Calculus"
		// and then show a Calculus box.
		expect(screen.getAllByText('Learn')).toHaveLength(1)
		expect(screen.getByRole('link', { name: 'Learn' })).toHaveAttribute('href', '/learn')
		expect(askAssistantMock).toHaveBeenCalledWith({ message: 'Where do I start?', history: [] })
	})

	it('sends earlier turns as history', async () => {
		getConfigMock.mockResolvedValue({ ai_enabled: true })
		askAssistantMock.mockResolvedValue({ reply: 'Yes.', segments: [] })
		render(<AssistantDock />)

		await screen.findByRole('button', { name: 'Open assistant' })
		openPanel()

		const input = await screen.findByLabelText('Your question')
		fireEvent.change(input, { target: { value: 'What is Mathua?' } })
		fireEvent.click(screen.getByRole('button', { name: 'Send' }))
		await screen.findByText('Yes.')

		fireEvent.change(input, { target: { value: 'And where do I start?' } })
		fireEvent.click(screen.getByRole('button', { name: 'Send' }))

		await waitFor(() => expect(askAssistantMock).toHaveBeenCalledTimes(2))
		expect(askAssistantMock.mock.calls[1][0]).toEqual({
			message: 'And where do I start?',
			history: [
				{ role: 'user', content: 'What is Mathua?' },
				{ role: 'assistant', content: 'Yes.' },
			],
		})
	})

	it('names a spent quota differently from an outage', async () => {
		getConfigMock.mockResolvedValue({ ai_enabled: true })
		askAssistantMock.mockRejectedValue(new Error('429 Too Many Requests'))
		render(<AssistantDock />)

		await screen.findByRole('button', { name: 'Open assistant' })
		openPanel()

		fireEvent.change(await screen.findByLabelText('Your question'), { target: { value: 'hi' } })
		fireEvent.click(screen.getByRole('button', { name: 'Send' }))

		expect(await screen.findByText(/today’s limit/i)).toBeInTheDocument()
	})
})
