import { ref } from "vue";
import { defineStore } from "pinia";
import { useAppStore } from "./app";
import { useCcStore } from "./chatting";

// This store talks directly to two external services (unlike every other
// store here, which only ever calls core-chatting's own Connect-Web API):
// - NoCloud's own REST/gateway API (app.conf.api), to list "bots"-type instances
// - ai-bot-manager's REST API, to link/unlink the core_chatting channel
// Both reuse the same iframe bearer token (app.conf.token), since all three
// services validate JWTs signed with the same shared NoCloud signing key.

export interface BotsInstance {
  uuid: string;
  title: string;
  status: string;
  botId: string; // ai-bot-manager Bot id, from instance.data.bot_uuid
}

export interface BotChannelLink {
  botId: string;
  botName: string;
  channelId: string;
  accountUuid: string; // NoCloud/core-chatting account acting as the bot's identity
  customName: string;
  skipReview: boolean;
}

export const useBotChannelStore = defineStore("bot_channel", () => {
  const app = useAppStore();
  const chatting = useCcStore();

  const instances = ref<BotsInstance[]>([]);
  const links = ref<BotChannelLink[]>([]);
  const isLoading = ref(false);

  function nocloudBase(): string {
    return (chatting?.baseUrl || "").replace(/\/$/, "");
  }

  function aiBotManagerBase(): string {
    const override = (import.meta as any).env?.VITE_AI_BOT_MANAGER_API as
      | string
      | undefined;
    const base = (override || chatting.baseUrl || "").replace(/\/$/, "");
    return `${base}/agents/api`;
  }

  function authHeaders(): HeadersInit {
    return {
      Authorization: `Bearer ${app.conf?.token}`,
      "Content-Type": "application/json",
    };
  }

  async function throwIfNotOk(res: Response, action: string) {
    if (!res.ok) {
      const text = await res.text().catch(() => "");
      throw new Error(`${action} failed (${res.status}): ${text}`);
    }
  }

  async function fetchInstances() {
    const res = await fetch(
      `${nocloudBase()}/nocloud.instances.InstancesService/List`,
      {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify({ filters: { type: ["bots"] } }),
      },
    );
    await throwIfNotOk(res, "Failed to fetch instances");
    const data = await res.json();

    // InstancesService/List returns a map keyed by index (not an array, and
    // not the pool/instancesGroups/instances nesting ServicesService/List
    // uses): { "0": { instance: {...}, type, service, sp, account, namespace }, ... }
    const raw = data.pool ?? data;
    const result: BotsInstance[] = [];
    for (const entry of Object.values(raw) as any[]) {
      const inst = entry?.instance;
      if (!inst) continue;
      const entryType = entry.type ?? inst.type;
      if (entryType !== "bots" || inst.status === "DEL") continue;
      const botId = inst.data?.bot_uuid;
      if (!botId) continue;
      result.push({
        uuid: inst.uuid,
        title: inst.title,
        status: inst.status,
        botId,
      });
    }
    instances.value = result;
  }

  async function fetchLinks() {
    const res = await fetch(`${aiBotManagerBase()}/get_bots`, {
      headers: authHeaders(),
    });
    await throwIfNotOk(res, "Failed to fetch bots");
    const data = await res.json();

    const result: BotChannelLink[] = [];
    for (const bot of data.bots || []) {
      for (const channel of bot.channels || []) {
        if (channel.type !== "core_chatting") continue;
        result.push({
          botId: bot.id,
          botName: bot.name,
          channelId: channel.id,
          accountUuid: channel.data?.bot_uuid || "",
          customName: channel.metadata?.custom_name || "",
          skipReview: !!channel.data?.skip_review,
        });
      }
    }
    links.value = result;
  }

  async function fetchAll() {
    isLoading.value = true;
    try {
      await Promise.all([fetchInstances(), fetchLinks()]);
    } finally {
      isLoading.value = false;
    }
  }

  async function addLink(
    botId: string,
    accountUuid: string,
    customName: string,
    skipReview: boolean,
  ) {
    const res = await fetch(`${aiBotManagerBase()}/add_channel`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify({
        bot: botId,
        type: "core_chatting",
        data: { bot_uuid: accountUuid, skip_review: skipReview },
        custom_name: customName,
      }),
    });
    await throwIfNotOk(res, "Failed to add channel");
  }

  async function removeLink(botId: string, channelId: string) {
    const res = await fetch(`${aiBotManagerBase()}/delete_channel`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify({ bot: botId, channel: channelId }),
    });
    await throwIfNotOk(res, "Failed to delete channel");
  }

  // Edit = delete the old channel then add the new one (ai-bot-manager has
  // no update_channel endpoint); throws before removing anything if the add
  // payload is invalid is not possible client-side, so callers should confirm
  async function editLink(
    link: BotChannelLink,
    accountUuid: string,
    customName: string,
    skipReview: boolean,
  ) {
    await removeLink(link.botId, link.channelId);
    await addLink(link.botId, accountUuid, customName, skipReview);
  }

  return {
    instances,
    links,
    isLoading,
    fetchAll,
    addLink,
    removeLink,
    editLink,
  };
});
