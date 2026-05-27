export const LANGUAGE_CODES = ["en", "es"] as const;
export type LanguageCode = typeof LANGUAGE_CODES[number];

export type TextValues = Record<LanguageCode, string>