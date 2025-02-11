<template>
  <div>
    <n-space align="center" justify="space-between" style="margin-bottom: 15px">
      <n-input
        style="min-width: 50vw"
        clearable
        placeholder="Start typing..."
        v-model:value="searchParam"
      />
      <n-button v-if="isEdit" @click="isAddTemplateActive = true">
        Add
      </n-button>
    </n-space>

    <n-space vertical justify="start" :wrap-item="false">
      <n-list>
        <n-scrollbar
          ref="scrollbar"
          style="max-height: calc(100dvh - 200px); margin-bottom: auto"
        >
          <n-list-item
            style="padding: 10px 0px; margin: 0px 10px"
            v-for="template in filtredTemplates"
            :key="template.name"
          >
            <n-space justify="space-between">
              <n-h6>{{ template.name }}</n-h6>

              <n-space style="margin-right: 20px">
                <template v-if="isEdit">
                  <n-button @click="startEditTemplate(template.name)" circle>
                    <template #icon>
                      <n-icon>
                        <pencil-icon />
                      </n-icon>
                    </template>
                  </n-button>

                  <n-button
                    @click="deleteTemplate(template.name)"
                    :loading="isSaveLoading"
                    circle
                  >
                    <template #icon>
                      <n-icon>
                        <delete-icon />
                      </n-icon>
                    </template>
                  </n-button>
                </template>

                <template v-else>
                  <n-button @click="emit('select', template.content)">
                    Select

                    <template #icon>
                      <n-icon>
                        <copy-icon />
                      </n-icon>
                    </template>
                  </n-button>
                </template>
              </n-space>
            </n-space>

            <!-- @vue-ignore -->
            <message-content
              :message="{ content: template.content } as any"
              :attach-files="[]"
            />
          </n-list-item>
        </n-scrollbar>
      </n-list>
    </n-space>
  </div>

  <n-modal v-model:show="isAddTemplateActive">
    <n-card
      style="width: 600px"
      title="Add new template"
      :bordered="false"
      size="huge"
      role="dialog"
      aria-modal="true"
    >
      <n-space vertical>
        <n-input placeholder="Template name" v-model:value="newTemplate.name" />

        <n-mention
          ref="input"
          size="small"
          type="textarea"
          style="width: 100%"
          placeholder="Type your message"
          v-model:value="newTemplate.content"
          :autosize="{ minRows: 5, maxRows: 15 }"
          :options="mentionsOptions"
          :prefix="isMentionVisible ? '@' : ' @'"
          @input="onInput"
        />
      </n-space>

      <template #footer>
        <n-space justify="end">
          <n-button @click="saveTemplate" :loading="isSaveLoading"
            >Save</n-button
          >
        </n-space>
      </template>
    </n-card>
  </n-modal>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent, ref, toRefs } from "vue";
import {
  NButton,
  NInput,
  NSpace,
  NH6,
  NIcon,
  NScrollbar,
  NList,
  NListItem,
  NMention,
  NModal,
  NCard,
  MentionOption,
  useNotification,
} from "naive-ui";
import messageContent from "../../components/chats/message/message_content.vue";
import { Defaults, Department } from "../../connect/cc/cc_pb";
import {
  MessageTemplate,
  MetricWithKey,
  useDefaultsStore,
} from "../../store/defaults";
import { useUsersStore } from "../../store/users";
import { storeToRefs } from "pinia";

const deleteIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CloseOutline")
);
const pencilIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/PencilOutline")
);
const copyIcon = defineAsyncComponent(
  () => import("@vicons/ionicons5/CopyOutline")
);

interface TemplatesProps {
  admins: string[];
  departments: Department[];
  gateways: string[];
  metrics: MetricWithKey[];
  templates: MessageTemplate[];
  isEdit: boolean;
}

const props = defineProps<TemplatesProps>();
const { admins, templates, departments, gateways, metrics, isEdit } =
  toRefs(props);

const emit = defineEmits(["refresh", "select"]);

const notification = useNotification();
const store = useDefaultsStore();
const usersStore = useUsersStore();
const defaultsStore = useDefaultsStore();

const { users } = storeToRefs(usersStore);

const searchParam = ref("");

const newTemplates = ref<MessageTemplate[]>([...templates.value]);

const newTemplate = ref<MessageTemplate>({
  content: "",
  name: "New template",
});
const oldTemplateName = ref("");
const isAddTemplateActive = ref(false);
const isMentionVisible = ref(false);
const isSaveLoading = ref(false);

const filtredTemplates = computed(() => {
  return newTemplates.value.filter(
    (template) =>
      !searchParam.value ||
      template.name.toLowerCase().includes(searchParam.value.toLowerCase()) ||
      template.content.toLowerCase().includes(searchParam.value.toLowerCase())
  );
});

const mentionsOptions = computed(() => {
  const result: MentionOption[] = [];

  admins.value.forEach((admin) => {
    const user = users.value.get(admin)!;
    result.push({
      label: user.title ?? admin,
      value: user.title?.replace(" ", "_") ?? admin,
    });
  });

  return result;
});

function onInput({ target }: { target: HTMLTextAreaElement }) {
  if (target.selectionStart !== target.selectionEnd) return;
  const symbol = target.value[target.selectionEnd - 1];
  const prevSymbol = target.value[target.selectionEnd - 2];

  if (symbol === "@" && prevSymbol === " ") {
    isMentionVisible.value = true;
  } else {
    isMentionVisible.value = false;
  }
}

async function saveTemplates() {
  try {
    isSaveLoading.value = true;

    const data = new Defaults({
      admins: admins.value,
      departments: departments.value,
      gateways: gateways.value,
      metrics: metrics.value.reduce(
        (result, metric) => ({
          ...result,
          [metric.key]: metric,
        }),
        {}
      ),
      bot: defaultsStore.bot,
      templates: newTemplates.value.reduce(
        (result, template) => ({
          ...result,
          [template.name]: template.content,
        }),
        {}
      ),
    });

    await store.update_defaults(data);

    notification.success({ title: "Done", duration: 1000 });

    emit("refresh");
  } catch (error) {
    notification.error({ title: (error as Error).message ?? error });
  } finally {
    setTimeout(() => {
      isSaveLoading.value = false;
    }, 500);
  }
}

function deleteTemplate(name: string) {
  newTemplates.value = newTemplates.value.filter((t) => t.name !== name);
  saveTemplates();
}

function saveTemplate() {
  if (oldTemplateName.value) {
    newTemplates.value = newTemplates.value.filter(
      (t) => t.name !== oldTemplateName.value
    );
  }

  newTemplates.value.push(newTemplate.value);

  newTemplate.value = {
    content: "",
    name: "New template",
  };
  oldTemplateName.value = "";

  saveTemplates();
  isAddTemplateActive.value = false;
}

function startEditTemplate(name: string) {
  const template = newTemplates.value.find((t) => t.name === name);
  if (!template) {
    return;
  }

  oldTemplateName.value = template.name;
  newTemplate.value = { ...template };
  isAddTemplateActive.value = true;
}
</script>
