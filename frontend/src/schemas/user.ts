import * as z from "zod";

export const UserSchema = z.object({
  id: z.uuid(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
  email: z.email(),
  language: z.enum(["en", "es"])
}).nullable();

export type User = z.infer<typeof UserSchema>;

