import { ref } from "vue";
import { defineStore } from "pinia";
import {
  aiBotManagerBase,
  authHeaders,
  nocloudBase,
  throwIfNotOk,
} from "./aiBotManager";

// Talks directly to the NoCloud REST gateway (list "bots" instances) and to
// ai-bot-manager's REST API (link/unlink the core_chatting channel). See
// ./aiBotManager for how the base URL and bearer token are shared.

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
  const instances = ref<BotsInstance[]>([]);
  const links = ref<BotChannelLink[]>([]);
  const isLoading = ref(false);

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

  // Edit = delete the old channel then re-add it (ai-bot-manager has no
  // update_channel endpoint). Not atomic: if the re-add fails the channel is
  // left removed, so the UI surfaces the error and the operator can retry.
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
