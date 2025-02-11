<template>
  <div>
    <n-space style="padding: 6px 24px" align="center">
      <n-h3>Update config</n-h3>
    </n-space>

    <n-space vertical justify="start" :wrap-item="false">
      <div v-for="(option, key) in options" :key="option.key">
        <n-space
          v-if="option.headers || option.isInput"
          align="center"
          :wrap-item="false"
        >
          {{ option.key }}
          <n-button
            ghost
            text
            size="small"
            :type="inputs[key]?.editMode ? 'error' : 'success'"
            @click="option.onClick"
          >
            <template #icon>
              <delete-icon v-if="inputs[key]?.editMode" />
              <plus-icon v-else />
            </template>
          </n-button>

          <n-space
            v-if="option.isInput && inputs[key].editMode"
            style="flex-wrap: nowrap"
            :wrap-item="false"
          >
            <n-input
              v-model:value="(inputs[key].value as string)"
              size="small"
            />
            <!-- @vue-ignore -->
            <n-button ghost text type="success" @click="option.onSave">
              <template #icon> <save-icon /> </template>
            </n-button>
          </n-space>
        </n-space>
        <n-text v-else>{{ option.key }}</n-text>

        <n-table v-if="option.headers" :single-line="false">
          <thead>
            <tr>
              <th v-for="header of option.headers">{{ header.title }}</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            <!-- @vue-ignore -->
            <tr v-for="(item, i) of config[key]">
              <td v-for="header of option.headers" :style="header.style">
                <n-space v-if="header.isEditable" :wrap-item="false">
                  <n-tag
                    round
                    closable
                    type="info"
                    v-for="(value, optionIndex) of item[header.value]"
                    :key="i"
                    @close="header.onClose && header.onClose(i, optionIndex)"
                  >
                    {{ value.key }} ({{ value.value }})
                  </n-tag>

                  <n-space
                    align="center"
                    style="flex-wrap: nowrap"
                    :wrap-item="false"
                  >
                    <template v-if="inputs[key].index === i && header.onSave">
                      <n-input
                        v-model:value="inputs[key].key"
                        size="small"
                        placeholder="Key"
                      />
                      <n-input-number
                        v-model:value="(inputs[key].value as number)"
                        size="small"
                        placeholder="Value"
                      />
                      <n-button
                        ghost
                        text
                        type="success"
                        @click="header.onSave(i)"
                      >
                        <template #icon> <save-icon /> </template>
                      </n-button>
                    </template>

                    <n-button
                      ghost
                      text
                      size="small"
                      :type="inputs[key]?.index === i ? 'error' : 'success'"
                      @click="header.onClick && header.onClick(i)"
                    >
                      <template #icon>
                        <delete-icon v-if="inputs[key]?.index === i" />
                        <plus-icon v-else />
                      </template>
                    </n-button>
                  </n-space>
                </n-space>

                <n-select
                  multiple
                  filterable
                  v-else-if="Array.isArray(item[header.value])"
                  v-model:value="item[header.value]"
                  value-field="uuid"
                  label-field="title"
                  :options="header.options"
                  style="min-width: 200px"
                />

                <n-switch
                  v-else-if="header.isBool"
                  v-model:value="item[header.value]"
                />
                <n-input v-else clearable v-model:value="item[header.value]" />
              </td>
              <td>
                <n-space justify="center" :wrap-item="false">
                  <n-button
                    ghost
                    text
                    type="error"
                    size="large"
                    @click="option.onClose && option.onClose(i)"
                  >
                    <template #icon> <delete-icon /> </template>
                  </n-button>
                </n-space>
              </td>
            </tr>
          </tbody>
        </n-table>

        <member-select-pagination
          v-else-if="option.isUsers"
          v-model:value="newAdmins"
        />

        <n-select
          v-else-if="option.items"
          multiple
          filterable
          v-model:value="(config[key as keyof configType] as Value)"
          value-field="uuid"
          label-field="title"
          :options="option.items"
          style="min-width: 200px"
        />

        <n-space v-else :wrap-item="false">
          <n-tag
            round
            closable
            type="info"
            v-for="(value, i) of config[key as keyof configType]"
            :key="i"
            @close="option.onClose && option.onClose(+i)"
          >
            {{ value }}
          </n-tag>
        </n-space>
      </div>

      <div>
        <n-space>
          <n-text>Bot</n-text>
        </n-space>

        <n-space>
          <n-switch v-model:value="config.bot.enable">
            <template #checked> Bot active </template>
            <template #unchecked> Bot disabled </template>
          </n-switch>

          <n-switch v-model:value="config.bot.review">
            <template #checked> Review by default </template>
            <template #unchecked> No review by default </template>
          </n-switch>
        </n-space>

        <n-space style="margin-top: 20px">
          <n-text>Promt</n-text>
        </n-space>
        <n-input
          v-model:value="config.bot.prompt"
          type="textarea"
          autosize
          placeholder="Bot Promt"
        />

        <n-space style="margin-top: 20px">
          <n-text>Custom values</n-text>
        </n-space>

        <n-table :bordered="false" :single-line="false">
          <thead>
            <tr>
              <th>Key</th>
              <th>Value</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(value, index) in botCustomValues">
              <td>
                <n-input v-model:value="value.key" />
              </td>
              <td>
                <n-input v-model:value="value.value" />
              </td>
              <td>
                <n-button
                  ghost
                  @click="deleteNewBotCustomValue(index)"
                  type="warning"
                  >Delete</n-button
                >
              </td>
            </tr>

            <tr>
              <td></td>
              <td></td>
              <td>
                <n-space justify="end">
                  <n-button ghost @click="addNewBotCustomValue" type="success"
                    >Add</n-button
                  >
                </n-space>
              </td>
            </tr>
          </tbody>
        </n-table>
      </div>

      <n-space justify="end" style="margin-bottom: 10px">
        <n-button :loading="isEditLoading" ghost type="info" @click="submit">
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
  NTag,
  useNotification,
} from "naive-ui";
import {
  computed,
  defineAsyncComponent,
  onMounted,
  reactive,
  toRefs,
  watch,
} from "vue";
import { Value } from "naive-ui/es/select/src/interface";
import {
  Bot,
  Defaults,
  Department,
  Metric,
  Option,
  User,
} from "../../connect/cc/cc_pb";
import { ref } from "vue";
import {
  MessageTemplate,
  MetricWithKey,
  useDefaultsStore,
} from "../../store/defaults.ts";
import MemberSelectPagination from "../users/member_select_pagination.vue";

const saveIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/SaveOutline")
);
const plusIcon = defineAsyncComponent(() => import("@vicons/ionicons5/Add"));
const deleteIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CloseOutline")
);

interface inputsType {
  [key: string]: {
    value: string | number;
    editMode?: boolean;
    index?: number;
    key?: string;
  };
}

interface optionsType {
  [key: string]: {
    key: string;
    headers?: [
      {
        title: string;
        value: string;
        options?: any[];
        isBool?: boolean;
        isEditable: boolean;
        style?: string;

        onClick?: (metricIndex: number) => void;
        onSave?: (i: number) => void;
        onClose?: (metricIndex: number, i: number) => void;
      }
    ];
    items?: any[];
    isInput?: boolean;
    isUsers?: boolean;

    onClick?: () => void;
    onSave?: (i?: number) => void;
    onClose?: (i: number) => void;
  };
}

interface configType {
  admins: User[];
  departments: Department[];
  gateways: string[];
  metrics: MetricWithKey[];
  templates: MessageTemplate[];
  bot: Bot;
}

const props = defineProps<configType>();
const { admins, departments, gateways, metrics, templates, bot } =
  toRefs(props);

const emit = defineEmits(["refresh"]);

const defaultsStore = useDefaultsStore();
const notification = useNotification();

const botCustomValues = ref<{ value: string; key: string }[]>();
const newAdmins = ref<string[]>([]);
const isEditLoading = ref(false);
const inputs = reactive<inputsType>({
  gateways: { editMode: false, value: "" },
  metrics: { index: -1, key: "", value: 0 },
});

const config = reactive<configType>({
  admins: [],
  departments: [],
  gateways: [],
  metrics: [],
  templates: [],
  bot: Bot.fromJson({}),
});

onMounted(() => {
  setDefaultConfig();
});

const departmentsOptions = computed(() => ({
  key: "Departments",
  headers: [
    { title: "Key", value: "key" },
    { title: "Title", value: "title" },
    { title: "Description", value: "description" },
    {
      title: "Admins",
      value: "admins",
      options: admins.value.map((admin) => {
        const { email } = (admin?.data as any) ?? {};

        return {
          ...admin,
          title: `${admin?.title} ${email ? `(${email})` : ""}`,
        };
      }),
    },
    { title: "WHMCS id", value: "whmcsId", style: "width: 50px" },
    { title: "Public", value: "public", isBool: true },
  ],
  onClick() {
    config.departments.push(
      new Department({
        key: "",
        title: "",
        description: "",
        whmcsId: "",
        admins: [],
        public: true,
      })
    );
  },
  onClose(i: number) {
    config.departments.splice(i, 1);
  },
}));

const gatewaysOptions = computed(() => ({
  key: "Gateways",
  isInput: true,
  onClick() {
    inputs.gateways.editMode = !inputs.gateways.editMode;
    inputs.gateways.value = "";
  },
  onSave() {
    config.gateways.push(inputs.gateways.value as string);
    inputs.gateways.editMode = false;
    inputs.gateways.value = "";
  },
  onClose(i: number) {
    config.gateways.splice(i, 1);
  },
}));

const metricsOptions = computed(() => ({
  key: "Metrics",
  headers: [
    { title: "Key", value: "key" },
    { title: "Title", value: "title" },
    {
      title: "Options",
      value: "options",
      isEditable: true,

      onClick(metricIndex: number) {
        if (inputs.metrics.index === metricIndex) {
          inputs.metrics = { index: -1, key: "", value: 0 };
        } else {
          inputs.metrics.index = metricIndex;
        }
      },
      onSave(i: number) {
        config.metrics[i].options.push(
          new Option({
            key: inputs.metrics.key,
            value: inputs.metrics.value as number,
          })
        );
        inputs.metrics = { index: -1, key: "", value: 0 };
      },
      onClose(metricIndex: number, i: number) {
        config.metrics[metricIndex].options.splice(i, 1);
      },
    },
  ],
  onClick() {
    const metric = new Metric({ options: [], title: "" }) as MetricWithKey;

    metric.key = "";
    config.metrics.push(metric);
  },
  onClose(i: number) {
    config.metrics.splice(i, 1);
  },
}));

const adminsOptions = computed(() => ({
  key: "Admins",
  isUsers: true,
}));

/* @ts-ignore */
const options = computed<optionsType>(() => ({
  admins: adminsOptions.value,
  departments: departmentsOptions.value,
  gateways: gatewaysOptions.value,
  metrics: metricsOptions.value,
}));

async function submit() {
  try {
    isEditLoading.value = true;
    await defaultsStore.update_defaults(
      new Defaults({
        ...config,
        admins: newAdmins.value,
        bot: {
          ...config.bot,
          values: botCustomValues.value?.reduce<any>((acc, curr) => {
            acc[curr.key] = curr.value;
            return acc;
          }, {}),
        },
        metrics: config.metrics.reduce(
          (result, metric) => ({
            ...result,
            [metric.key]: metric,
          }),
          {}
        ),
        templates: config.templates.reduce(
          (result, template) => ({
            ...result,
            [template.name]: template.content,
          }),
          {}
        ),
      })
    );

    notification.success({ title: "Done", duration: 1000 });

    emit("refresh");
  } catch (error) {
    notification.error({ title: (error as Error).message ?? error });
  } finally {
    isEditLoading.value = false;
  }
}

function setDefaultConfig() {
  config.metrics = metrics.value;
  config.admins = admins.value;
  config.gateways = gateways.value;
  config.departments = departments.value;
  config.templates = templates.value;
  config.bot = bot.value || Bot.fromJson({});
}

function addNewBotCustomValue() {
  botCustomValues.value?.push({ key: "", value: "" });
}

function deleteNewBotCustomValue(index: number) {
  botCustomValues.value = botCustomValues.value?.filter(
    (_, current) => current !== index
  );
}

watch(
  () => config.bot.values,
  () => {
    botCustomValues.value = Object.keys(config.bot.values).map((key) => ({
      key,
      value: config.bot.values[key],
    }));
  }
);

newAdmins.value = admins.value.map((a) => a.uuid);

watch(admins, () => {
  newAdmins.value = admins.value.map((a) => a.uuid);
});
</script>

<script lang="ts">
export default {
  name: "settings-config",
};
</script>
