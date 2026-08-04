export const LANGUAGE_CODES = ["en", "es"] as const;
export type LanguageCode = typeof LANGUAGE_CODES[number];

export function isLanguageCode(value: string): value is LanguageCode {
  const normalized = value.trim().toLowerCase();
  const langCodesArray = LANGUAGE_CODES as readonly string[];
  if (langCodesArray.includes(normalized)) {
    return true
  } else {
    return false
  }
}

export type TextValues = Record<LanguageCode, string>