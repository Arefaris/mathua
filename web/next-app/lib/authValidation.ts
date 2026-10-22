import type { TFunction } from './i18n'

// Client-side checks on the sign-up and login fields. They exist to say what is wrong before a
// request is made; the server enforces the same rules and stays the authority. Each message
// takes the translate function from the component, because the text is shown to the learner.

export type FieldErrors = { name?: string; username?: string; email?: string; password?: string }

// Mirrors internal/auth password policy — keep both in sync:
// min 8 chars (Go: len bytes), max 72 bytes (bcrypt limit), a Unicode digit,
// and a special = neither Unicode letter nor digit (Go: unicode.IsLetter).
export function passwordIssues(pw: string, t: TFunction): string[] {
  const issues: string[] = []
  if (pw.length < 8) issues.push(t('account.login.issueMin'))
  if (new TextEncoder().encode(pw).length > 72) issues.push(t('account.login.issueMax'))
  if (!/\p{Nd}/u.test(pw)) issues.push(t('account.login.issueNumber'))
  if (!/[^\p{L}\p{Nd}]/u.test(pw)) issues.push(t('account.login.issueSpecial'))
  return issues
}

/** "Password needs at least 8 characters, a number" */
export function passwordNeeds(issues: string[], t: TFunction): string {
  return t('account.login.passwordNeeds', { issues: issues.join(', ') })
}

// Mirrors internal/auth ValidateUsername on the normalized (lowercased,
// trimmed) form — keep both in sync.
function usernameIssues(username: string, t: TFunction): string | null {
  const u = username.trim().toLowerCase()
  if (u.length < 3 || [...u].length > 20) return t('account.login.usernameLength')
  if (!/^[a-z0-9_.]+$/.test(u)) return t('account.login.usernameChars')
  return null
}

// Mirrors internal/auth ValidateEmail on the trimmed form — keep both in
// sync. Empty is handled by validateFields (required on signup).
function emailIssues(email: string, t: TFunction): string | null {
  const e = email.trim()
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(e)) return t('account.login.emailInvalid')
  return null
}

export function validateFields(
  tab: 'login' | 'signup',
  values: { name: string; username: string; email: string; password: string },
  t: TFunction,
): FieldErrors {
  const errors: FieldErrors = {}
  if (tab === 'signup' && !values.name.trim()) errors.name = t('account.login.nameRequired')
  if (!values.username.trim()) {
    errors.username = t('account.login.usernameRequired')
  } else if (tab === 'signup') {
    const issue = usernameIssues(values.username, t)
    if (issue) errors.username = issue
  }
  if (tab === 'signup') {
    if (!values.email.trim()) {
      errors.email = t('account.login.emailRequired')
    } else {
      const issue = emailIssues(values.email, t)
      if (issue) errors.email = issue
    }
  }
  if (!values.password) {
    // aislop-ignore-next-line security/hardcoded-secret -- validation message string, not a credential
    errors.password = t('account.login.passwordRequired')
  } else if (tab === 'signup') {
    const issues = passwordIssues(values.password, t)
    // aislop-ignore-next-line security/hardcoded-secret -- validation message string, not a credential
    if (issues.length > 0) errors.password = passwordNeeds(issues, t)
  }
  return errors
}
