// Otus's three modes, named in one place. The install-wide switch (Settings ->
// Bot) and a single chat's override (the chat's Bot settings dialog) are the
// same three values; describing them two different ways in two dialogs is how
// operators end up believing they are two different products.
//
// The key is read from two different bags depending on which control wrote it:
// cc.Defaults.Bot.values for the install switch, and the chat's own bot_state
// for the per-chat override, which wins where it is set.
export const OTUS_MODE_KEY = "otus.mode";

export interface OtusModeOption {
  value: string;
  label: string;
  hint: string;
}

export const OTUS_MODES: OtusModeOption[] = [
  {
    value: "copilot_mode",
    label: "copilot_mode",
    hint: "Answers the operator in the copilot lane. Writes to the client only when the operator asked it to.",
  },
  {
    value: "self_mode",
    label: "self_mode",
    hint: "May write to the client. Still asks in the copilot lane when a person has to decide.",
  },
  {
    value: "god_mode",
    label: "god_mode",
    hint: "Same as self_mode, and may change the ticket status.",
  },
];

// A chat with no override of its own stores an empty string; the picker has no
// rung for it, so opening the dialog preselects the install-wide mode instead.
export const OTUS_CHAT_MODE_INHERIT = "";

export function otusModeLabel(value: string | undefined): string {
  return OTUS_MODES.find((option) => option.value === value)?.label ?? "not set";
}
