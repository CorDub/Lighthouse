import * as z from "zod";

export const DefaultsSchema = z.object({
  lang: z.enum(["en", "es"]).default("en")
})

export type Defaults = z.infer<typeof DefaultsSchema>;