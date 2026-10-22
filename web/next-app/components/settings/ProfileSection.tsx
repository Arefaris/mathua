'use client'

import { useState } from 'react'
import { Input } from '@/components/ui/input'
import { updateSettings, updateProfileName, uploadAvatar, deleteAvatar, avatarImageUrl, requestEmailVerification, type UserSettings } from '../../lib/api'
import { getUserInfo, setUserInfo } from '../../lib/auth'
import { DICEBEAR_STYLES, dicebearUrl, randomDicebear, type DicebearPick } from '../../lib/dicebear'
import { useT } from '../../lib/i18n'
import { Dices } from '../icons'
import Loading from '../Loading'
import Avatar from '../Avatar'
import RichText from '../RichText'

// Profile helpers: every save strips the legacy avatar_preset key (quiet
// migration to the DiceBear model — presets are no longer offered).
function stripPreset(s: UserSettings): UserSettings {
  const next = { ...s }
  delete next.avatar_preset
  return next
}

/**
 * The Profile block of /settings: display name, recovery email, and the profile picture.
 *
 * It was inline in the page, which had grown past the file-length limit. The settings object
 * stays the page's: this block reads it and hands every change back through `onSettings`, so
 * the page's other controls (timer, pause) keep saving against the same copy.
 */
export default function ProfileSection({ settings, onSettings }: { settings: UserSettings; onSettings: (next: UserSettings) => void }) {
  const t = useT()
  // Read at mount, not in an effect: this block only renders after the page has loaded in the
  // browser, and an effect would paint one frame with an empty name field.
  const [displayName, setDisplayName] = useState(() => getUserInfo()?.name ?? '')
  const [nameBusy, setNameBusy] = useState(false)
  const [nameMsg, setNameMsg] = useState('')
  const [recoveryEmail, setRecoveryEmail] = useState(() => getUserInfo()?.email ?? '')
  const [photoBusy, setPhotoBusy] = useState(false)
  const [photoMsg, setPhotoMsg] = useState('')
  const [photoVersion, setPhotoVersion] = useState<number | null>(null)
  const [pendingDice, setPendingDice] = useState<{ style: DicebearPick['style']; seed: string } | null>(null)
  const [imageBusy, setImageBusy] = useState(false)
  const [imageMsg, setImageMsg] = useState('')
  const [verifyMsg, setVerifyMsg] = useState('')
  const [verifyBusy, setVerifyBusy] = useState(false)

  const gallerySeed = (typeof window !== 'undefined' ? getUserInfo()?.student_id : undefined) ?? 'guest'

  // One effective seed everywhere: a staged surprise re-seeds the whole
  // grid (new character, different outfits); otherwise the saved pick's
  // seed; otherwise the stable per-user seed. After Save, pending clears
  // and the grid falls back to the just-saved seed — no visual jump.
  const gridSeed = pendingDice?.seed ?? settings.avatar_dicebear?.seed ?? gallerySeed

  const handleSaveName = async () => {
    const name = displayName.trim()
    if (!name) {
      setNameMsg(t('account.settings.nameEmpty'))
      return
    }
    setNameBusy(true)
    setNameMsg('')
    try {
      const res = await updateProfileName(name, recoveryEmail.trim())
      const info = getUserInfo()
      if (info) setUserInfo({ ...info, name: res.name, email: res.email || info.email })
      setDisplayName(res.name)
      setNameMsg(t('account.settings.profileSaved'))
    } catch (e) {
      // Surface friendly server messages (e.g. "email already in use");
      // fall back to the generic notice for network failures.
      const serverMessage = (e as { serverMessage?: string })?.serverMessage
      setNameMsg(serverMessage || t('account.settings.saveFailed'))
    } finally {
      setNameBusy(false)
    }
  }

  // Gallery picks and Surprise only stage a preview — Save image persists.
  const handlePickStyle = (style: DicebearPick['style']) => {
    const seed = pendingDice?.seed ?? settings.avatar_dicebear?.seed ?? gallerySeed
    setPendingDice({ style, seed })
    setImageMsg('')
  }

  const handleSurprise = () => {
    const pick = randomDicebear()
    setPendingDice({ style: pick.style, seed: pick.seed })
    setImageMsg('')
  }

  const pendingDiffers =
    pendingDice !== null &&
    (pendingDice.style !== settings.avatar_dicebear?.style ||
      pendingDice.seed !== settings.avatar_dicebear?.seed ||
      settings.avatar_custom === true)

  const handleSaveImage = async () => {
    if (!pendingDice) return
    setImageBusy(true)
    setImageMsg('')
    try {
      const next = stripPreset({ ...settings, avatar_custom: false, avatar_dicebear: { style: pendingDice.style, seed: pendingDice.seed } })
      onSettings(next)
      await updateSettings(next)
      if (settings.avatar_custom) await deleteAvatar().catch(() => {})
      setPhotoVersion(null)
      setPendingDice(null)
      setImageMsg(t('account.settings.imageSaved'))
      // Broadcast so the Header picture refreshes in-tab (same object —
      // setUserInfo always dispatches auth-changed).
      const savedInfo = getUserInfo()
      if (savedInfo) setUserInfo({ ...savedInfo })
    } catch {
      setImageMsg(t('account.settings.saveFailed'))
    } finally {
      setImageBusy(false)
    }
  }

  const handleUpload = async (file: File | undefined) => {
    if (!file) return
    if (file.size > 512 * 1024) {
      setPhotoMsg(t('account.settings.photoTooBig'))
      return
    }
    setPhotoBusy(true)
    setPhotoMsg('')
    try {
      await uploadAvatar(file)
      const next = stripPreset({ ...settings, avatar_custom: true, avatar_version: Date.now() })
      onSettings(next)
      setPendingDice(null)
      setPhotoVersion(Date.now())
      setPhotoMsg(t('account.settings.photoUpdated'))
      const uploadedInfo = getUserInfo()
      if (uploadedInfo) setUserInfo({ ...uploadedInfo })
    } catch {
      setPhotoMsg(t('account.settings.uploadFailed'))
    } finally {
      setPhotoBusy(false)
    }
  }

  const handleRemovePhoto = async () => {
    setPhotoBusy(true)
    try {
      await deleteAvatar()
      const next = stripPreset({ ...settings, avatar_custom: false, avatar_version: Date.now() })
      onSettings(next)
      setPendingDice(null)
      setPhotoVersion(null)
      const removedInfo = getUserInfo()
      if (removedInfo) setUserInfo({ ...removedInfo })
    } catch { console.error('deleteAvatar failed') }
    finally { setPhotoBusy(false) }
  }

  const handleVerifyEmail = async () => {
    setVerifyBusy(true)
    setVerifyMsg('')
    try {
      await requestEmailVerification()
      setVerifyMsg(t('account.settings.verifySent'))
      const info = getUserInfo()
      if (info) setUserInfo({ ...info })
    } catch (e: unknown) {
      setVerifyMsg(e instanceof Error ? e.message : t('account.login.requestFailed'))
    } finally {
      setVerifyBusy(false)
    }
  }

  const currentPhotoUrl = settings.avatar_custom ? avatarImageUrl(photoVersion ?? undefined) : undefined
  const currentDice = settings.avatar_dicebear

  const userEmail = typeof window !== 'undefined' ? getUserInfo()?.email ?? '' : ''
  const userEmailVerified = typeof window !== 'undefined' ? getUserInfo()?.email_verified ?? false : false

  return (
    <div className="border-t border-mathua-border pt-6 min-w-0">
      <span className="font-mono text-sm text-mathua-primary">{t('account.settings.profileTitle')}</span>
      <p className="text-mathua-muted text-xs mt-1">
        {t('account.settings.profileHint')}
      </p>
      <label htmlFor="display-name" className="font-mono text-[11px] uppercase text-mathua-muted mt-4 block">{t('account.settings.displayName')}</label>
      <form onSubmit={e => { e.preventDefault(); handleSaveName() }} className="flex flex-col sm:flex-row gap-2 mt-2">
        <Input
          id="display-name"
          type="text"
          value={displayName}
          maxLength={50}
          onChange={(e) => setDisplayName(e.target.value)}
          placeholder={t('account.settings.namePlaceholder')}
          enterKeyHint="done"
          className="flex-1 h-14"
        />
        <button
          type="submit"
          disabled={nameBusy}
          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none px-4 h-10 text-xs font-mono disabled:opacity-50 shrink-0 inline-flex items-center gap-2"
        >
          {nameBusy ? (<><Loading inline size={11} /> {t('common.action.saving')}</>) : t('account.settings.saveProfile')}
        </button>
      </form>
      {nameMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{nameMsg}</p>}

      <label htmlFor="recovery-email" className="font-mono text-[11px] uppercase text-mathua-muted mt-4 block">{t('account.settings.recoveryEmail')}</label>
      <p className="text-mathua-muted text-xs mt-1">{t('account.settings.recoveryHint')}</p>
      <div className="flex flex-col sm:flex-row gap-2 mt-2">
        <Input
          id="recovery-email"
          type="email"
          value={recoveryEmail}
          onChange={(e) => setRecoveryEmail(e.target.value)}
          placeholder="you@example.com"
          autoComplete="email"
          className="flex-1 h-14"
        />
      </div>
      {userEmail !== '' && !userEmailVerified && (
        <div className="mt-2 flex flex-wrap items-center gap-2">
          <span className="font-mono text-[11px] text-mathua-muted">{t('account.settings.emailUnverified')}</span>
          <button
            type="button"
            onClick={handleVerifyEmail}
            disabled={verifyBusy}
            className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-3 h-8 text-[11px] font-mono disabled:opacity-50"
          >
            {verifyBusy ? t('account.login.sending') : t('account.settings.verifyEmail')}
          </button>
        </div>
      )}
      {verifyMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{verifyMsg}</p>}

      <div className="flex items-center gap-4 mt-6">
        <Avatar
          seed={gridSeed}
          name={displayName || t('account.settings.you')}
          size={56}
          url={pendingDice ? dicebearUrl(pendingDice.style, pendingDice.seed) : currentPhotoUrl ?? (currentDice ? dicebearUrl(currentDice.style as DicebearPick['style'], currentDice.seed) : undefined)}
        />
        <div className="min-w-0">
          <p className="font-mono text-[11px] uppercase text-mathua-muted">{t('account.settings.currentPicture')}</p>
          <p className="text-mathua-muted text-xs mt-1 break-words">
            {pendingDice
              ? t('account.settings.pictureStaged')
              : settings.avatar_custom ? t('account.settings.pictureUploaded') : currentDice ? t('account.settings.pictureDicebear') : t('account.settings.pictureAuto')}
          </p>
        </div>
      </div>

      <p className="font-mono text-[11px] uppercase text-mathua-muted mt-6 mb-2">{t('account.settings.pickCharacter')}</p>
      <div className="grid grid-cols-4 sm:grid-cols-6 gap-2">
        {DICEBEAR_STYLES.map((style) => {
          const shown = pendingDice ?? currentDice
          const selected = !settings.avatar_custom && shown?.style === style
          return (
            <button
              type="button"
              key={style}
              onClick={() => handlePickStyle(style)}
              title={style}
              aria-label={t('account.settings.pickStyle', { style })}
              className={`p-1 border rounded-none transition-colors bg-mathua-code ${selected ? 'border-mathua-blue' : 'border-transparent hover:border-mathua-border'}`}
            >
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img src={dicebearUrl(style, gridSeed)} alt={t('account.settings.styleAlt', { style })} width={64} height={64} loading="lazy" className="w-full h-auto" />
            </button>
          )
        })}
      </div>
      <div className="flex flex-wrap gap-2 mt-3">
        <button
          type="button"
          onClick={handleSurprise}
          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none px-4 h-9 text-xs font-mono inline-flex items-center gap-2"
        >
          <Dices className="size-4" aria-hidden />
          {t('account.settings.surprise')}
        </button>
        <button
          type="button"
          onClick={handleSaveImage}
          disabled={!pendingDiffers || imageBusy}
          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none px-4 h-9 text-xs font-mono disabled:opacity-50 inline-flex items-center gap-2"
        >
          {imageBusy ? (<><Loading inline size={11} /> {t('common.action.saving')}</>) : t('account.settings.saveImage')}
        </button>
        <label className="border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue rounded-none px-4 h-9 text-xs font-mono inline-flex items-center gap-2 cursor-pointer focus-within:border-mathua-blue focus-within:text-mathua-blue">
          {photoBusy ? (<><Loading inline size={11} /> {t('account.settings.uploading')}</>) : t('account.settings.uploadPhoto')}
          <input
            type="file"
            accept="image/png,image/jpeg,image/gif,image/webp"
            aria-label={t('account.settings.uploadAria')}
            aria-describedby="photo-limits"
            className="hidden"
            disabled={photoBusy}
            onChange={(e) => { handleUpload(e.target.files?.[0]); e.target.value = '' }}
          />
        </label>
        {settings.avatar_custom && (
          <button
            type="button"
            onClick={handleRemovePhoto}
            disabled={photoBusy}
            className="border border-mathua-border text-mathua-secondary hover:border-mathua-red hover:text-mathua-red rounded-none px-4 h-9 text-xs font-mono disabled:opacity-50"
          >
            {t('account.settings.removePhoto')}
          </button>
        )}
      </div>
      {photoMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{photoMsg}</p>}
      {imageMsg && <p className="font-mono text-[11px] text-mathua-secondary mt-2">{imageMsg}</p>}
      <p id="photo-limits" className="text-mathua-muted text-[11px] mt-3">
        <RichText
          text={t('account.settings.dicebearCredit')}
          slots={{ dicebear: <a href="https://www.dicebear.com" target="_blank" rel="noreferrer" className="text-mathua-blue hover:text-mathua-blue-hover">DiceBear</a> }}
        />
      </p>
    </div>
  )
}
