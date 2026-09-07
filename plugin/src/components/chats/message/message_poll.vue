<template>
  <div class="poll" :style="{ gridColumn: appStore.isMobile ? '1 / 3' : '2 / 3' }">
    <div v-if="poll.question" class="poll__question">{{ poll.question }}</div>

    <div class="poll__options">
      <!-- One answer: the option itself is the button. Several: tick them and
           send, or every tick would be an answer of its own. -->
      <template v-if="poll.multiple">
        <n-checkbox
          v-for="option of poll.options"
          :key="option.id"
          :checked="picked.includes(option.id)"
          :disabled="disabled"
          @update:checked="(on: boolean) => toggle(option.id, on)"
        >
          {{ option.label }}
        </n-checkbox>
      </template>
      <template v-else>
        <n-button
          v-for="option of poll.options"
          :key="option.id"
          class="poll__option"
          :type="answered.includes(option.id) ? 'primary' : 'default'"
          :secondary="!answered.includes(option.id)"
          :disabled="disabled"
          :loading="busy === option.id"
          @click="answer([option.id])"
        >
          {{ option.label }}
        </n-button>
      </template>
    </div>

    <n-button
      v-if="poll.multiple"
      class="poll__send"
      size="small"
      type="primary"
      :disabled="disabled || !picked.length || sameAsAnswered"
      :loading="busy === MULTI"
      @click="answer(picked)"
    >
      {{ answered.length ? "Change answer" : "Answer" }}
    </n-button>

    <!-- No way to take an answer back: picking another one replaces it, which
         is the only thing anybody wanted from that. -->
    <div class="poll__foot">
      <!-- The operator's side of a poll is the one line that matters: what the
           customer picked. Everything else about it they can already read. -->
      <n-text v-if="others.length" depth="3">{{ others.join(" · ") }}</n-text>
      <n-text v-else-if="answered.length" depth="3">
        Your answer is saved, you can change it
      </n-text>
      <n-text v-else-if="poll.closed" depth="3">This poll is closed</n-text>
      <n-text v-else depth="3">Pick an answer</n-text>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, toRefs } from "vue";
import { NButton, NCheckbox, NText, useNotification } from "naive-ui";
import { createPromiseClient } from "@connectrpc/connect";
import { createGrpcWebTransport } from "@connectrpc/connect-web";
import { Message, VoteRequest } from "../../../connect/cc/cc_pb";
import { MessagesAPI } from "../../../connect/cc/cc_connect";
import { useCcStore } from "../../../store/chatting";
import { useAppStore } from "../../../store/app";

const props = defineProps<{ message: Message }>();
const { message } = toRefs(props);

const store = useCcStore();
const appStore = useAppStore();
const notification = useNotification();

const MULTI = "__multi__";
const busy = ref("");

// The client is built here rather than taken from the chat store, which keeps
// its own. Reason: pinia 2.1.6 — the version this image is built with — loses
// the inferred action types of that store as soon as one more member is added
// to it, and the build fails on an unrelated line. Not worth reshaping a store
// this size over one call.
//
// ponytail: fold this back into the store once pinia is upgraded (the lockfile
// bump to 3.x is already in the tree).
const transport = createGrpcWebTransport({
  baseUrl: import.meta.env.VITE_API_URL || "/",
  useBinaryFormat: true,
  interceptors: [
    (next) => async (req) => {
      req.header.set("Authorization", `Bearer ${appStore.conf?.token}`);
      return next(req);
    },
  ],
});
const messages = createPromiseClient(MessagesAPI, transport);

const poll = computed(() => message.value.poll!);

/** answered is what I picked, as stored on the server. */
const answered = computed(() => poll.value.votes[store.me.uuid]?.options ?? []);

/** picked is what is ticked but not sent yet; it starts from my answer. */
const picked = ref<string[]>([...answered.value]);

const sameAsAnswered = computed(
  () =>
    picked.value.length === answered.value.length &&
    picked.value.every((id) => answered.value.includes(id))
);

const disabled = computed(() => poll.value.closed || busy.value !== "");

/**
 * Everybody else's answers, by option: in a ticket this is the customer's
 * answer as the operator sees it. Names are not resolved — a ticket has one
 * customer, and "Too expensive" is the whole of what an operator needs.
 */
const others = computed(() =>
  Object.entries(poll.value.votes)
    .filter(([account]) => account !== store.me.uuid)
    .flatMap(([, vote]) =>
      vote.options.map(
        (id) => poll.value.options.find((o) => o.id === id)?.label ?? id
      )
    )
);

function toggle(id: string, on: boolean) {
  picked.value = on
    ? [...picked.value, id]
    : picked.value.filter((existing) => existing !== id);
}

async function answer(options: string[]) {
  busy.value = poll.value.multiple ? MULTI : options[0];
  try {
    // One call: the ticket service records the answer and keeps the single
    // message that states it in words — posting it, editing it when the
    // answer changes, removing it when the answer is retracted. Doing that
    // here is what turned four clicks into four replies in the chat.
    await messages.vote(new VoteRequest({ message: message.value.uuid, options }));
    picked.value = [...options];
  } catch (e: any) {
    notification.error({ content: e.message ?? String(e), duration: 5000 });
  } finally {
    busy.value = "";
  }
}
</script>

<style scoped>
.poll {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
  padding: 10px 12px;
  border: 1px solid var(--n-border-color, #8884);
  border-radius: 8px;
  max-width: 420px;
}
.poll__question {
  font-weight: 600;
}
.poll__options {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
/* Answers read as a list, so they are left-aligned rather than centred like
   an ordinary button. */
.poll__option {
  justify-content: flex-start;
  text-align: left;
  height: auto;
  min-height: 30px;
  padding: 6px 10px;
  white-space: normal;
}
.poll__send {
  align-self: flex-start;
}
.poll__foot {
  font-size: 12px;
}
</style>
