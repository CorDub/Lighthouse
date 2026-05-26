import type { LanguageCode, TextValues } from "../Text"

export function getText(textValue: TextValues, langCode: LanguageCode): string {
  return textValue[langCode]
} 