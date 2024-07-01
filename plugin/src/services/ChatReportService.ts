import writeXlsxFile from "write-excel-file";
import { Chat, Status, User } from "../connect/cc/cc_pb";
import { useCcStore } from "../store/chatting";

interface ResponsibleReportItem {
  name: string;
  all: number;
  opened: number;
  closed: number;
  activity: number;
  request_mediana: number;
  request_periods: number[];
}

export default class ChatReportService {
  // - в идеале - среднее время ответа (оценить как он реагирует - быстро или медленно)

  private static schema = [
    {
      column: "Имя",
      type: String,
      value: (item: ResponsibleReportItem) => item.name,
      width: 20,
    },
    {
      column: "Всего",
      type: Number,
      value: (item: ResponsibleReportItem) => item.all,
      width: 15,
    },
    {
      column: "Закрыто",
      type: Number,
      value: (item: ResponsibleReportItem) => item.closed,
      width: 15,
    },
    {
      column: "Открыто",
      type: Number,
      value: (item: ResponsibleReportItem) => item.opened,
      width: 15,
    },
    {
      column: "Сколько сообщений написал",
      type: Number,
      value: (item: ResponsibleReportItem) => item.activity,
      width: 35,
    },
    {
      column: "Среднее время ответа в минутах",
      type: String,
      value: (item: ResponsibleReportItem) =>
        item.request_mediana ? item.request_mediana.toFixed(2).toString() : "-",
      width: 35,
    },
  ];

  static async generate(chats: Chat[], admins: User[]) {
    const csStore = useCcStore();

    const responsibleReportItems: Map<string, ResponsibleReportItem> =
      new Map();

    const promises = chats.map(async (chat) => {
      if (chat.responsible) {
        const responsible = responsibleReportItems.get(chat.responsible) || {
          activity: 0,
          all: 0,
          closed: 0,
          opened: 0,
          name: admins.find((a) => a.uuid === chat.responsible)?.title || "",
          request_mediana: 0,
          request_periods: [],
        };

        responsible.all += 1;
        if (chat.status === Status.CLOSE) {
          responsible.closed += 1;
        } else {
          responsible.opened += 1;
        }

        responsibleReportItems.set(chat.responsible, responsible);
      }

      const messagesData = await csStore.get_messages(chat, false);

      let isFirstAdminResponseExist = false;
      messagesData.messages.forEach((message) => {
        if (!admins.find((admin) => admin.uuid === message.sender)) {
          return;
        }

        const responsible = responsibleReportItems.get(message.sender) || {
          activity: 0,
          all: 0,
          closed: 0,
          opened: 0,
          name: admins.find((a) => a.uuid === message.sender)?.title || "",
          request_mediana: 0,
          request_periods: [],
        };

        responsible.activity += 1;

        if (!isFirstAdminResponseExist) {
          isFirstAdminResponseExist = true;

          const firstUserMessageSent = Number(messagesData.messages[0].sent);
          const firstAdminMessageSent = Number(message.sent);

          responsible.request_periods.push(
            (firstAdminMessageSent - firstUserMessageSent) / 1000 / 60
          );
        }

        responsibleReportItems.set(message.sender, responsible);
      });
    });

    await Promise.all(promises);

    await writeXlsxFile(
      [...responsibleReportItems.values()].map((item) => ({
        ...item,
        request_mediana:
          item.request_periods.reduce((acc, item) => acc + item, 0) /
          item.request_periods.length,
      })),
      {
        schema: this.schema,
        fileName: `report_${new Date().toISOString().split('T')[0]}.xlsx`,
      }
    );
  }
}
