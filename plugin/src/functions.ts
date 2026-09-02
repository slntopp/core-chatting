import { NotificationApi } from "naive-ui";
import { Kind, Message, Status } from "./connect/cc/cc_pb";

export function addToClipboard(text: string, notification?: NotificationApi) {
  if (navigator?.clipboard) {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        if (!notification) {
          alert("Text copied");
          return;
        }
        notification.success({ content: "Text copied", duration: 3000 });
      })
      .catch((res) => {
        console.error(res);
      });
  } else {
    if (!notification) {
      alert("Clipboard is not supported");
      return;
    }
    notification.error({ content: "Clipboard is not supported" });
  }
}

export function getImageUrl(name: string) {
  const icons = import.meta.glob(`/**/*.png`);
  let key = Object.keys(icons).find((key) => key.includes(`/${name}`)) ?? "";

  if (!key.includes("cc.ui")) key = `/cc.ui${key}`;
  return key.replace("/dist", "").replace("/public", "");
}

export function getRelativeTime(
  timestamp: number,
  now: number,
  isLifetime?: boolean
) {
  const timeDifference = (now - timestamp) / 1000;
  const minutesDifference = Math.floor(timeDifference / 60);

  if (minutesDifference >= 4320 && !isLifetime) {
    return new Date(timestamp).toLocaleDateString();
  } else if (minutesDifference >= 1440) {
    const daysDifference = Math.floor(minutesDifference / 1440);

    return `${daysDifference} days${isLifetime ? "" : " ago"}`;
  } else if (minutesDifference >= 60) {
    const hoursDifference = Math.floor(minutesDifference / 60);

    return `${hoursDifference} hours${isLifetime ? "" : " ago"}`;
  } else if (minutesDifference > 0) {
    return `${minutesDifference} minutes${isLifetime ? "" : " ago"}`;
  } else {
    return "just now";
  }
}

export function getStatusColor(status: Status) {
  switch (status) {
    case Status.NEW:
      return "#5084ff";
    case Status.ON_HOLD:
      return "#00dbff";
    case Status.OPEN:
      return "#1ea01e";
    case Status.IN_PROGRESS:
      return "#00ffaa";
    case Status.CUSTOMER_REPLY:
      return "#ffcc55";
    case Status.RESOLVE:
    case Status.ANSWERED:
      return "#ff8300";
    case Status.CLOSE:
      return "#e23535";
    // Outreach tickets: deliberately grey. They are ours, not a customer's
    // request, and nothing about them is urgent.
    case Status.ONBOARDING:
      return "#8f8f8f";
    default:
      return undefined;
  }
}

export function getStatusItems() {
  const allowedStatuses = [0, 1, 8, 5, 4, 7, 3, 9];

  return allowedStatuses.map((status) => ({
    label: Status[status].replace("_", " "),
    value: status,
  }));
}

export const debounce = (callback: Function, wait: number) => {
  let timeoutId: number;
  return (...args: any) => {
    window.clearTimeout(timeoutId);
    timeoutId = window.setTimeout(() => {
      callback(...args);
    }, wait);
  };
};

export function cleanObject(object: any) {
  Object.entries(object).forEach(([k, v]) => {
    if (v && typeof v === "object") {
      cleanObject(v);
    }
    if (
      (v && typeof v === "object" && !Object.keys(v).length) ||
      v === null ||
      v === undefined
    ) {
      if (Array.isArray(object)) {
        object.splice(k as any, 1);
      } else {
        delete object[k];
      }
    }
  });
  return object;
}

/**
 * Splits a chat's messages into the two panes of the operator UI.
 *
 * `customer` is the conversation as the customer will see it: their messages,
 * approved replies, admin-only notes, and at most ONE pending draft - the
 * newest, the one the Approve button acts on. `copilot` is the conversation
 * with the AI assistant plus every draft that a later revision has already
 * superseded. Those earlier drafts are working history, not messages anyone
 * received, and leaving them in the customer thread makes it read as a pile of
 * duplicate answers.
 *
 * Admin notes belong in the main thread. They are what operators write to each
 * other about this chat, in its own timeline - the customer cannot see them
 * either way, and filing them under the copilot mixed two unrelated things:
 * a note left for a colleague kept waking the bot, and the bot's answers
 * buried the notes.
 *
 * Lives here rather than in either component so the two can never disagree
 * about which draft is still live and show it twice, or not at all.
 */
/**
 * Kinds the customer never receives: notes operators leave for each other, and
 * the operator's conversation with the AI. Mirrors cc.OperatorOnlyKinds on the
 * server - anything deciding "is this part of the customer conversation" must
 * ask this and not compare against one kind, which is how a second hidden lane
 * would leak the first time someone forgot a branch.
 */
export function isOperatorOnly(kind: Kind) {
  return kind === Kind.ADMIN_ONLY || kind === Kind.COPILOT;
}

export function splitChatMessages(all: Message[]) {
  const sorted = [...all].sort((a, b) => Number(a.sent - b.sent));

  let livePending = "";
  for (let i = sorted.length - 1; i >= 0; i--) {
    if (sorted[i].kind !== Kind.COPILOT && sorted[i].underReview) {
      livePending = sorted[i].uuid;
      break;
    }
  }

  const customer: Message[] = [];
  const copilot: Message[] = [];
  for (const m of sorted) {
    if (m.kind === Kind.COPILOT) copilot.push(m);
    else if (m.underReview && m.uuid !== livePending) copilot.push(m);
    else customer.push(m);
  }
  return { customer, copilot };
}
