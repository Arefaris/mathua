'use client'

import * as Dialog from '@radix-ui/react-dialog'
import Link from 'next/link'
import type { EligibilityPrereq } from '../lib/api'
import type { ReturnEntry } from '../lib/learnReturn'
import { prereqHref, prereqStateWord } from './PrerequisitePanel'
import { useT } from '../lib/i18n'

/**
 * The prerequisite-path dialog for one locked or upcoming curriculum topic.
 *
 * Opened from the curriculum browser's "View prerequisite path" row action. It
 * renders the selected target's own title plus that target's actual readiness
 * response, never a hardcoded list, with every prerequisite linking to its
 * reference page. Reading here awards nothing; practice stays gated in /learn.
 */

export type PrereqDialogStatus = 'loading' | 'ready' | 'empty' | 'error'

interface Props {
  open: boolean
  targetTitle: string
  prerequisites: EligibilityPrereq[]
  status: PrereqDialogStatus
  error?: string
  returnStack?: ReturnEntry[]
  onClose: () => void
  onRetry?: () => void
}

export default function PrerequisitePathDialog({
  open,
  targetTitle,
  prerequisites,
  status,
  error,
  returnStack = [],
  onClose,
  onRetry,
}: Props) {
  const t = useT()
  return (
    <Dialog.Root
      open={open}
      onOpenChange={(next) => {
        if (!next) onClose()
      }}
    >
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-50 bg-black/50" />
        <Dialog.Content
          aria-describedby="prereq-path-desc"
          className="fixed left-1/2 top-1/2 z-50 max-h-[85vh] w-[calc(100vw-2rem)] max-w-[520px] -translate-x-1/2 -translate-y-1/2 overflow-y-auto border border-mathua-border bg-mathua-surface p-5"
        >
          <Dialog.Title className="font-mono text-xs uppercase tracking-wider text-mathua-primary">
            {t('learn.prereqDialog.title')}
          </Dialog.Title>
          <p className="mt-1 font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere]">
            {targetTitle}
          </p>
          <Dialog.Description
            id="prereq-path-desc"
            className="mt-2 font-mono text-[11px] text-mathua-secondary break-words [overflow-wrap:anywhere]"
          >
            {t('learn.prereqDialog.description')}
          </Dialog.Description>

          <h3 className="mt-4 font-mono text-[11px] uppercase tracking-wider text-mathua-muted">
            {t('learn.prereqDialog.heading')}
          </h3>

          {status === 'loading' && (
            <p className="mt-2 font-mono text-[11px] text-mathua-muted">{t('learn.prereqDialog.loading')}</p>
          )}

          {status === 'error' && (
            <div className="mt-2">
              <p
                role="alert"
                className="font-mono text-[11px] text-mathua-secondary break-words [overflow-wrap:anywhere]"
              >
                {error || t('learn.prereqDialog.error')}
              </p>
              <div className="mt-3 flex flex-wrap gap-2">
                {onRetry && (
                  <button
                    type="button"
                    onClick={onRetry}
                    className="border border-mathua-blue px-4 py-2 font-mono text-[11px] text-mathua-blue hover:bg-mathua-blue-faint min-h-[36px]"
                  >
                    {t('learn.prereqDialog.retry')}
                  </button>
                )}
              </div>
            </div>
          )}

          {status === 'empty' && (
            <p className="mt-2 font-mono text-[11px] text-mathua-muted">{t('learn.prereqDialog.empty')}</p>
          )}

          {status === 'ready' && (
            <ul className="mt-3 space-y-2">
              {prerequisites.map((p) => (
                <li
                  key={p.conceptId}
                  className="flex flex-col sm:flex-row sm:items-center gap-2 min-w-0"
                >
                  <span className="min-w-0 flex-1">
                    <span className="block font-mono text-[11px] text-mathua-primary break-words [overflow-wrap:anywhere]">
                      {p.title}
                    </span>
                    <span className="block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
                      {prereqStateWord(p)} · {t('learn.prereq.required')}
                    </span>
                  </span>
                  <Link
                    href={prereqHref(p, returnStack)}
                    className="w-full sm:w-auto shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap"
                  >
                    {t('learn.prereq.study')}
                  </Link>
                </li>
              ))}
            </ul>
          )}

          <div className="mt-4 flex justify-end">
            <Dialog.Close asChild>
              <button
                type="button"
                className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted hover:text-mathua-blue min-h-[36px] px-2"
              >
                {t('common.action.close')}
              </button>
            </Dialog.Close>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
