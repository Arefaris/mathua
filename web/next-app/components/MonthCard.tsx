'use client'

import { useMemo } from 'react'
import type { DailyActivity } from '../lib/api'
import { heatColor } from '../lib/heatColor'
import { formatDate, useLocale, useT, type Locale } from '../lib/i18n'

// Keys of the weekday headers, Sunday first to match Date#getDay. Only Monday, Wednesday and
// Friday are shown, as on the yearly grid, so the column labels do not crowd the cells.
export const DAY_KEYS = ['su', 'mo', 'tu', 'we', 'th', 'fr', 'sa'] as const
const SHOWN_DAY_INDEXES = new Set([1, 3, 5])

// The full month name in the learner's language, standalone and capitalized ("March", "Март").
// Intl rather than a list so a third language needs no new table.
export function monthName(month: number, locale: Locale): string {
  const name = new Date(2000, month, 1).toLocaleString(locale === 'en' ? 'en-US' : locale, { month: 'long' })
  return name.charAt(0).toUpperCase() + name.slice(1)
}

interface Props {
  year: number
  month: number
  data: DailyActivity[]
}

export default function MonthCard({ year, month, data }: Props) {
  const t = useT()
  const locale = useLocale()
  const monthData = useMemo(() => {
    const prefix = `${year}-${String(month + 1).padStart(2, '0')}`
    return data.filter((d) => d.date.startsWith(prefix))
  }, [data, year, month])

  const { calendar, stats } = useMemo(() => {
    const firstDay = new Date(year, month, 1)
    const lastDay = new Date(year, month + 1, 0)
    const daysInMonth = lastDay.getDate()
    const startDow = firstDay.getDay()

    const activityMap = new Map<string, { questions: number; correct: number; concepts: string[] }>()
    for (const d of data) {
      activityMap.set(d.date, { questions: d.questions, correct: d.correct, concepts: d.concepts })
    }

    const calendar: ({ date: number; questions: number; correct: number; level: number } | null)[][] = []
    let day = 1
    for (let week = 0; week < 6; week++) {
      const row: ({ date: number; questions: number; correct: number; level: number } | null)[] = []
      for (let dow = 0; dow < 7; dow++) {
        if ((week === 0 && dow < startDow) || day > daysInMonth) {
          row.push(null)
        } else {
          const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
          const act = activityMap.get(dateStr)
          const questions = act?.questions ?? 0
          const correct = act?.correct ?? 0
          row.push({ date: day, questions, correct, level: 0 })
          day++
        }
      }
      calendar.push(row)
    }

    const totalQ = monthData.reduce((sum, d) => sum + d.questions, 0)
    const totalC = monthData.reduce((sum, d) => sum + d.correct, 0)
    const allConcepts = new Set<string>()
    for (const d of monthData) {
      for (const cid of d.concepts) allConcepts.add(cid)
    }
    const accuracy = totalQ > 0 ? Math.round((totalC / totalQ) * 100) : 0

    let maxQ = 0
    for (const d of monthData) {
      if (d.questions > maxQ) maxQ = d.questions
    }

    for (const week of calendar) {
      for (const cell of week) {
        if (cell) {
          if (cell.questions === 0 || maxQ === 0) {
            cell.level = 0
          } else {
            const ratio = cell.questions / maxQ
            if (ratio <= 0.25) cell.level = 1
            else if (ratio <= 0.50) cell.level = 2
            else if (ratio <= 0.75) cell.level = 3
            else cell.level = 4
          }
        }
      }
    }

    return {
      calendar,
      stats: { questions: totalQ, correct: totalC, accuracy, concepts: allConcepts.size },
    }
  }, [data, year, month, monthData])

  const activeDays = monthData.filter((d) => d.questions > 0).length

  return (
    <div className="w-[272px] shrink-0 snap-start border-[0.5px] border-mathua-border overflow-hidden bg-transparent">
      <div className="px-3 py-2 min-w-0">
        <p className="min-w-0 truncate text-left text-[13px] font-normal font-mono text-mathua-primary">
          {monthName(month, locale)} {year}
          <span className="ml-2 text-[11px] text-mathua-muted">
            {t('common.count.days', { count: activeDays })}
          </span>
        </p>
      </div>

      <div className="px-3 pb-3 pt-0 border-t-[0.5px] border-mathua-border w-full min-w-0 overflow-hidden">
        <div className="flex flex-wrap gap-x-4 gap-y-1 py-2 font-mono text-[11px] text-mathua-secondary">
          <span>{t('common.count.questions', { count: stats.questions })}</span>
          <span className="text-mathua-muted">{t('goals.activity.correctPercent', { percent: stats.accuracy })}</span>
          <span className="text-mathua-muted">{t('goals.activity.topics', { count: stats.concepts })}</span>
        </div>

        <div className="flex justify-center w-full max-w-full min-w-0 overflow-hidden">
          <div className="inline-block max-w-full">
            <div className="flex gap-1 mb-1">
              {DAY_KEYS.map((dayKey, dayIndex) => (
                <div
                  key={dayKey}
                  className="w-5 h-5 flex items-center justify-center font-mono text-[10px] text-mathua-muted shrink-0"
                >
                  {SHOWN_DAY_INDEXES.has(dayIndex) ? t(`goals.activity.day.${dayKey}`) : ''}
                </div>
              ))}
            </div>

            {calendar.map((week, wi) => (
              <div key={wi} className="flex gap-1 mb-1">
                {week.map((cell, di) => (
                  <div
                    key={di}
                    title={
                      cell
                        ? t('goals.activity.cellTitle', {
                            date: formatDate(new Date(year, month, cell.date), locale, { month: 'short', day: 'numeric' }),
                            questions: cell.questions,
                            percent: cell.questions > 0 ? Math.round((cell.correct / cell.questions) * 100) : 0,
                          })
                        : ''
                    }
                    className="w-5 h-5 shrink-0 border-[0.5px]"
                    style={{
                      background: cell ? heatColor(cell.level) : 'transparent',
                      borderColor: cell ? 'var(--border)' : 'transparent',
                    }}
                  />
                ))}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
