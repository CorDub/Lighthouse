import * as z from "zod";

export const ConnectionSchema = z.object({
  id: z.uuid(),
  service: z.enum(["youtube", "instagram", "tiktok", "twitter", "twitch"]),
  channelId: z.string(),
  channelHandle: z.string(),
  scopes: z.string()
})

export type Connection = z.infer<typeof ConnectionSchema>;