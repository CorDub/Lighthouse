type TextValues = {
  "en": string,
  "es": string
}

export function getText(textValue: TextValues, langCode: "en" | "es"): string {
  return textValue[langCode]
} 