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
    label: "Hints",
    hint: "Answers the operator in the copilot lane. Writes to the client only when the operator asked it to.",
  },
  {
    value: "self_mode",
    label: "Writes itself",
    hint: "May write to the client. Still asks in the copilot lane when a person has to decide.",
  },
  {
    value: "god_mode",
    label: "Full",
    hint: 'Same as "Writes itself", and may change the ticket status.',
  },
];

// The per-chat control has a fourth rung the install switch cannot have: no
// override at all. It is the default, and an empty value is what clears one.
export const OTUS_CHAT_MODE_INHERIT = "";

export const OTUS_CHAT_MODES: OtusModeOption[] = [
  {
    value: OTUS_CHAT_MODE_INHERIT,
    label: "As in settings",
    hint: "This chat follows the install-wide switch on the Bot settings page.",
  },
  ...OTUS_MODES,
];

export function otusModeLabel(value: string | undefined): string {
  return OTUS_MODES.find((option) => option.value === value)?.label ?? "not set";
}
