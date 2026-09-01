<template>
  <div>
    <n-space style="padding: 6px 24px" align="center">
      <n-h3>Auto tickets</n-h3>
    </n-space>

    <n-space vertical style="padding: 0 24px; max-width: 720px">
      <n-text depth="3" style="display: block; margin-bottom: 12px">
        Creates a support ticket after the instance payment date plus the delay
        for that billing period. Placeholders: {CLIENT_NAME}, {INSTANCE},
        {PRODUCT}, {IPS}, {SERVICE_DETAILS}, {SERVICE}
      </n-text>

      <div class="bots_config_switch">
        <n-switch class="switch" v-model:value="autoTicket.enabled">
          <template #checked> Active </template>
          <template #unchecked> Disabled </template>
        </n-switch>
        <span> Enable automatic overdue tickets. </span>
      </div>

      <n-text>Default delay (hours)</n-text>
      <n-input-number
        v-model:value="autoTicket.defaultDelayHours"
        :min="0"
        placeholder="24"
      />
      <n-text depth="3">
        Used when the instance period is not listed below. Current default: 24
        hours.
      </n-text>

      <n-text>Delay by billing period</n-text>
      <n-table :single-line="false">
        <thead>
          <tr>
            <th>Period</th>
            <th>Delay (hours)</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, i) in autoTicket.delays" :key="i">
            <td>
              <n-select
                v-model:value="row.period"
                filterable
                tag
                :options="periodOptions"
                placeholder="Billing period"
              />
            </td>
            <td>
              <n-input-number v-model:value="row.hours" :min="0" />
            </td>
            <td>
              <n-button
                ghost
                type="error"
                @click="autoTicket.delays.splice(i, 1)"
              >
                Delete
              </n-button>
            </td>
          </tr>
        </tbody>
      </n-table>
      <n-button ghost type="success" @click="addDelayRule">Add period</n-button>

      <n-text>Department</n-text>
      <n-select
        v-model:value="autoTicket.department"
        clearable
        filterable
        placeholder="Department"
        label-field="title"
        value-field="key"
        :options="departmentSelectOptions"
      />

      <n-text>Assignees</n-text>
      <n-select
        v-model:value="autoTicket.admins"
        multiple
        clearable
        filterable
        placeholder="Extra assignees (in addition to department admins)"
        label-field="title"
        value-field="uuid"
        :options="adminSelectOptions"
      />

      <n-text>Responsible</n-text>
      <n-select
        v-model:value="autoTicket.responsible"
        clearable
        filterable
        placeholder="Responsible"
        label-field="title"
        value-field="uuid"
        :options="adminSelectOptions"
      />

      <n-text>Sender (WHMCS staff)</n-text>
      <n-select
        v-model:value="autoTicket.senderUuid"
        clearable
        filterable
        placeholder="Staff UUID used to open the ticket"
        label-field="title"
        value-field="uuid"
        :options="adminSelectOptions"
      />

      <n-text>Topic</n-text>
      <n-input v-model:value="autoTicket.topic" placeholder="Ticket topic" />

      <n-text>Message</n-text>
      <n-input
        v-model:value="autoTicket.message"
        type="textarea"
        :autosize="{ minRows: 8, maxRows: 16 }"
        placeholder="First message"
      />

      <n-space justify="end" style="margin: 10px 0 20px">
        <n-button :loading="isSaving" ghost type="info" @click="submit">
          Update
        </n-button>
      </n-space>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import {
  NButton,
  NInput,
  NInputNumber,
  NTable,
  NSelect,
  NSwitch,
  NSpace,
  NText,
  NH3,
  useNotification,
} from "naive-ui";
import { computed, reactive, ref, toRefs, watch } from "vue";
import { Bot, Defaults, Department, User } from "../../connect/cc/cc_pb";
import { MetricWithKey, useDefaultsStore } from "../../store/defaults.ts";

interface AutoTicketsProps {
  admins: User[];
  departments: Department[];
}

const props = defineProps<AutoTicketsProps>();
const { admins, departments } = toRefs(props);

const emit = defineEmits(["refresh"]);

const defaultsStore = useDefaultsStore();
const notification = useNotification();
const isSaving = ref(false);

const AUTO_TICKET_KEYS = {
  enabled: "auto_ticket.enabled",
  department: "auto_ticket.department",
  topic: "auto_ticket.topic",
  message: "auto_ticket.message",
  senderUuid: "auto_ticket.sender_uuid",
  admins: "auto_ticket.admins",
  responsible: "auto_ticket.responsible",
  delayHours: "auto_ticket.delay_hours",
  delays: "auto_ticket.delays",
};

const PERIOD_OPTIONS = [
  { label: "Hourly — 3600s", value: 3600 },
  { label: "Daily — 86400s", value: 86400 },
  { label: "Weekly — 604800s", value: 604800 },
  { label: "Monthly — 2592000s (30d)", value: 2592000 },
  { label: "Monthly — 2678400s (31d)", value: 2678400 },
  { label: "Yearly — 31536000s (365d)", value: 31536000 },
  { label: "Yearly — 32140800s (12×31d)", value: 32140800 },
];

const DEFAULT_AUTO_TICKET_TOPIC =
  "Уведомление об удалении услуги: {INSTANCE}";
const DEFAULT_AUTO_TICKET_MESSAGE = `Здравствуйте.

Уважаемый {CLIENT_NAME}, сообщаем, что оказание услуги: "{SERVICE_DETAILS}" приостановлено в связи с истечением срока оплаты.

Информация о выставленных счетах доступна в личном кабинете.
Обращаем внимание, что в случае неоплаты счета, размещенные данные будут удалены без возможности восстановления.
Если Вам нужна помощь, пожалуйста, свяжитесь с нами.

С уважением, служба поддержки.`;

function parseAutoTicket(values: Record<string, string> | undefined) {
  const src = values ?? {};
  const enabledRaw = (src[AUTO_TICKET_KEYS.enabled] ?? "true").toLowerCase();
  const delayHoursRaw = Number(src[AUTO_TICKET_KEYS.delayHours]);
  let delays: { period: number; hours: number }[] = [];
  try {
    const parsed = JSON.parse(src[AUTO_TICKET_KEYS.delays] || "[]");
    if (Array.isArray(parsed)) {
      delays = parsed
        .map((row) => ({
          period: Number(row?.period),
          hours: Number(row?.hours),
        }))
        .filter(
          (row) => Number.isFinite(row.period) && Number.isFinite(row.hours)
        );
    }
  } catch {
    delays = [];
  }
  return {
    enabled:
      enabledRaw !== "false" && enabledRaw !== "0" && enabledRaw !== "no",
    department: src[AUTO_TICKET_KEYS.department] ?? "",
    admins: (src[AUTO_TICKET_KEYS.admins] ?? "")
      .split(",")
      .map((v) => v.trim())
      .filter(Boolean),
    responsible: src[AUTO_TICKET_KEYS.responsible] || null,
    senderUuid: src[AUTO_TICKET_KEYS.senderUuid] || null,
    topic: src[AUTO_TICKET_KEYS.topic] || DEFAULT_AUTO_TICKET_TOPIC,
    message: src[AUTO_TICKET_KEYS.message] || DEFAULT_AUTO_TICKET_MESSAGE,
    defaultDelayHours:
      Number.isFinite(delayHoursRaw) && delayHoursRaw >= 0
        ? delayHoursRaw
        : 24,
    delays,
  };
}

function serializeAutoTicket(ticket: {
  enabled: boolean;
  department: string;
  admins: string[];
  responsible: string | null;
  senderUuid: string | null;
  topic: string;
  message: string;
  defaultDelayHours: number;
  delays: { period: number; hours: number }[];
}) {
  const delays = (ticket.delays ?? [])
    .map((row) => ({
      period: Number(row.period),
      hours: Number(row.hours),
    }))
    .filter(
      (row) =>
        Number.isFinite(row.period) &&
        row.period > 0 &&
        Number.isFinite(row.hours) &&
        row.hours >= 0
    );
  return {
    [AUTO_TICKET_KEYS.enabled]: ticket.enabled ? "true" : "false",
    [AUTO_TICKET_KEYS.department]: ticket.department ?? "",
    [AUTO_TICKET_KEYS.admins]: (ticket.admins ?? []).join(","),
    [AUTO_TICKET_KEYS.responsible]: ticket.responsible ?? "",
    [AUTO_TICKET_KEYS.senderUuid]: ticket.senderUuid ?? "",
    [AUTO_TICKET_KEYS.topic]: ticket.topic ?? "",
    [AUTO_TICKET_KEYS.message]: ticket.message ?? "",
    [AUTO_TICKET_KEYS.delayHours]: String(ticket.defaultDelayHours ?? 24),
    [AUTO_TICKET_KEYS.delays]: JSON.stringify(delays),
  };
}

const autoTicket = reactive(parseAutoTicket(defaultsStore.bot?.values));
const periodOptions = PERIOD_OPTIONS;

const departmentSelectOptions = computed(() =>
  (departments.value || []).map((dep) => ({
    key: dep.key,
    title: dep.title ? `${dep.title} (${dep.key})` : dep.key,
  }))
);

const adminSelectOptions = computed(() =>
  (admins.value || []).map((admin) => {
    const { email } = (admin?.data as any) ?? {};
    return {
      uuid: admin?.uuid,
      title: `${admin?.title ?? admin?.uuid ?? ""}${email ? ` (${email})` : ""}`,
    };
  })
);

function addDelayRule() {
  autoTicket.delays.push({ period: 2592000, hours: 24 });
}

async function submit() {
  isSaving.value = true;
  try {
    const bot = defaultsStore.bot;
    await defaultsStore.update_defaults(
      new Defaults({
        admins: defaultsStore.admins,
        departments: defaultsStore.departments,
        gateways: defaultsStore.gateways,
        metrics: defaultsStore.metrics.reduce(
          (result: Record<string, MetricWithKey>, metric) => ({
            ...result,
            [metric.key]: metric,
          }),
          {}
        ),
        templates: defaultsStore.templates.reduce(
          (result: Record<string, string>, template) => ({
            ...result,
            [template.name]: template.content,
          }),
          {}
        ),
        bot: new Bot({
          enable: bot?.enable ?? false,
          review: bot?.review ?? false,
          initiator: bot?.initiator ?? false,
          emergency: bot?.emergency ?? false,
          prompt: bot?.prompt ?? "",
          values: {
            ...(bot?.values ?? {}),
            ...serializeAutoTicket(autoTicket),
          },
        }),
      })
    );
    notification.success({ title: "Done", duration: 1000 });
    emit("refresh");
  } catch (error) {
    notification.error({ title: (error as Error).message ?? String(error) });
  } finally {
    isSaving.value = false;
  }
}

watch(
  () => defaultsStore.bot,
  (bot) => {
    Object.assign(autoTicket, parseAutoTicket(bot?.values));
  }
);
</script>

<script lang="ts">
export default {
  name: "settings-auto-tickets",
};
</script>

<style scoped>
.bots_config_switch {
  display: flex;
  margin: 5px 0px;
  align-items: center;
}

.bots_config_switch span {
  font-size: 1.2rem;
}

.bots_config_switch .switch {
  width: 200px;
  justify-content: normal;
}
</style>
