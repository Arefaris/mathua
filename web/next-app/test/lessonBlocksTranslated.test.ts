import { describe, it, expect } from 'vitest'
import { buildExampleBlocks, normSection } from '../lib/lessonBlocks'
import type { LessonKpsRes } from '../lib/api'

// The client never looks a step's section up in a lesson body: the server resolves it and sends
// the worked example as text. `section` is only used as part of the key that merges two steps
// pointing at the same text. So a translated lesson (Russian headings, Russian body) and an
// English one go through the same code, and these tests pin that down.

const russian: LessonKpsRes = {
  concept_id: 'geo.basic.angles',
  kps: [
    {
      label: 'Измерение «раствора» угла',
      section: 'Измерение «раствора» угла',
      subgoals: ['Сравните раствор с эталоном прямого угла'],
      worked_example: '## Измерение «раствора» угла\n\nПрямой угол равен ровно $90^\\circ$.',
    },
    {
      label: 'Тот же раздел под другим названием',
      section: 'Измерение «раствора» угла',
      subgoals: [],
      worked_example: '## Измерение «раствора» угла\n\nПрямой угол равен ровно $90^\\circ$.',
    },
    { label: 'Тренировка', section: '', subgoals: ['Назовите развёрнутый угол'], worked_example: '' },
  ],
}

describe('worked-example blocks of a translated lesson', () => {
  it('merges steps that point at the same translated section and keeps the rest', () => {
    const blocks = buildExampleBlocks(russian.concept_id, russian)
    expect(blocks.map(b => b.label)).toEqual(['Измерение «раствора» угла', 'Тренировка'])
    expect(blocks[0].also).toEqual(['Тот же раздел под другим названием'])
    expect(blocks[0].body).toContain('Прямой угол равен ровно')
    expect(blocks[1].body).toBe('')
  })

  it('treats a translated section like an English one', () => {
    expect(normSection('Решение  \\(x\\)  и $y$')).toBe('Решение x и y')
    expect(normSection('Solving  \\(x\\)  and $y$')).toBe('Solving x and y')
  })
})
