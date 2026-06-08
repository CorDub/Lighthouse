import { createSignal } from "solid-js";

export type AgencyHomeModule = "empty" | "createReport"

export const [displayed, setDisplayed] = createSignal<AgencyHomeModule>("empty")