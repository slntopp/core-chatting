# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: cc/cc.proto
# Protobuf Python Version: 5.29.3
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    3,
    '',
    'cc/cc.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import struct_pb2 as google_dot_protobuf_dot_struct__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0b\x63\x63/cc.proto\x12\x02\x63\x63\x1a\x1cgoogle/protobuf/struct.proto\"\x07\n\x05\x45mpty\"\x81\x02\n\x08\x43hatMeta\x12\x16\n\x06unread\x18\x01 \x01(\x05R\x06unread\x12\x30\n\rfirst_message\x18\x02 \x01(\x0b\x32\x0b.cc.MessageR\x0c\x66irstMessage\x12.\n\x0clast_message\x18\x03 \x01(\x0b\x32\x0b.cc.MessageR\x0blastMessage\x12*\n\x04\x64\x61ta\x18\x04 \x03(\x0b\x32\x16.cc.ChatMeta.DataEntryR\x04\x64\x61ta\x1aO\n\tDataEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12,\n\x05value\x18\x02 \x01(\x0b\x32\x16.google.protobuf.ValueR\x05value:\x02\x38\x01\"\x8c\x04\n\x04\x43hat\x12\x12\n\x04uuid\x18\x01 \x01(\tR\x04uuid\x12\x14\n\x05users\x18\x02 \x03(\tR\x05users\x12\x16\n\x06\x61\x64mins\x18\x03 \x03(\tR\x06\x61\x64mins\x12\x14\n\x05owner\x18\x04 \x01(\tR\x05owner\x12\x1a\n\x08gateways\x18\x05 \x03(\tR\x08gateways\x12\x1c\n\x04role\x18\x06 \x01(\x0e\x32\x08.cc.RoleR\x04role\x12\x19\n\x05topic\x18\x07 \x01(\tH\x00R\x05topic\x88\x01\x01\x12%\n\x04meta\x18\x08 \x01(\x0b\x32\x0c.cc.ChatMetaH\x01R\x04meta\x88\x01\x01\x12\x18\n\x07\x63reated\x18\t \x01(\x03R\x07\x63reated\x12\"\n\x06status\x18\n \x01(\x0e\x32\n.cc.StatusR\x06status\x12\x1e\n\ndepartment\x18\x0b \x01(\tR\ndepartment\x12%\n\x0bresponsible\x18\x0c \x01(\tH\x02R\x0bresponsible\x88\x01\x01\x12\x33\n\tbot_state\x18\r \x03(\x0b\x32\x16.cc.Chat.BotStateEntryR\x08\x62otState\x1aS\n\rBotStateEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12,\n\x05value\x18\x02 \x01(\x0b\x32\x16.google.protobuf.ValueR\x05value:\x02\x38\x01\x42\x08\n\x06_topicB\x07\n\x05_metaB\x0e\n\x0c_responsible\"\'\n\x05\x43hats\x12\x1e\n\x05\x63hats\x18\x01 \x03(\x0b\x32\x08.cc.ChatR\x05\x63hats\"\x1d\n\x05Merge\x12\x14\n\x05\x63hats\x18\x01 \x03(\tR\x05\x63hats\"\xb1\x02\n\x10ListChatsRequest\x12\x17\n\x04page\x18\x01 \x01(\x04H\x00R\x04page\x88\x01\x01\x12\x19\n\x05limit\x18\x02 \x01(\x04H\x01R\x05limit\x88\x01\x01\x12\x19\n\x05\x66ield\x18\x03 \x01(\tH\x02R\x05\x66ield\x88\x01\x01\x12\x17\n\x04sort\x18\x04 \x01(\tH\x03R\x04sort\x88\x01\x01\x12;\n\x07\x66ilters\x18\x05 \x03(\x0b\x32!.cc.ListChatsRequest.FiltersEntryR\x07\x66ilters\x1aR\n\x0c\x46iltersEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12,\n\x05value\x18\x02 \x01(\x0b\x32\x16.google.protobuf.ValueR\x05value:\x02\x38\x01\x42\x07\n\x05_pageB\x08\n\x06_limitB\x08\n\x06_fieldB\x07\n\x05_sort\"\xc8\x02\n\x12SetBotStateRequest\x12\x12\n\x04\x63hat\x18\x01 \x01(\tR\x04\x63hat\x12\x37\n\x05state\x18\x02 \x03(\x0b\x32!.cc.SetBotStateRequest.StateEntryR\x05state\x12\x1f\n\x08\x64isabled\x18\x03 \x01(\x08H\x00R\x08\x64isabled\x88\x01\x01\x12!\n\tescalated\x18\x04 \x01(\x08H\x01R\tescalated\x88\x01\x01\x12$\n\x0bskip_review\x18\x05 \x01(\x08H\x02R\nskipReview\x88\x01\x01\x1aP\n\nStateEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12,\n\x05value\x18\x02 \x01(\x0b\x32\x16.google.protobuf.ValueR\x05value:\x02\x38\x01\x42\x0b\n\t_disabledB\x0c\n\n_escalatedB\x0e\n\x0c_skip_review\"G\n\x11ListChatsResponse\x12\x1c\n\x04pool\x18\x01 \x03(\x0b\x32\x08.cc.ChatR\x04pool\x12\x14\n\x05total\x18\x02 \x01(\x03R\x05total\"\xa5\x01\n\x11\x43ountChatsRequest\x12<\n\x07\x66ilters\x18\x01 \x03(\x0b\x32\".cc.CountChatsRequest.FiltersEntryR\x07\x66ilters\x1aR\n\x0c\x46iltersEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12,\n\x05value\x18\x02 \x01(\x0b\x32\x16.google.protobuf.ValueR\x05value:\x02\x38\x01\"\x93\x01\n\x12\x43ountChatsResponse\x12@\n\x08statuses\x18\x01 \x03(\x0b\x32$.cc.CountChatsResponse.StatusesEntryR\x08statuses\x1a;\n\rStatusesEntry\x12\x10\n\x03key\x18\x01 \x01(\x05R\x03key\x12\x14\n\x05value\x18\x02 \x01(\x03R\x05value:\x02\x38\x01\":\n\nAttachment\x12\x12\n\x04type\x18\x01 \x01(\tR\x04type\x12\x18\n\x07\x63ontent\x18\x02 \x01(\x0cR\x07\x63ontent\"\xd0\x03\n\x07Message\x12\x12\n\x04uuid\x18\x01 \x01(\tR\x04uuid\x12\x1c\n\x04kind\x18\x02 \x01(\x0e\x32\x08.cc.KindR\x04kind\x12\x16\n\x06sender\x18\x03 \x01(\tR\x06sender\x12\x18\n\x07\x63ontent\x18\x04 \x01(\tR\x07\x63ontent\x12 \n\x0b\x61ttachments\x18\x05 \x03(\tR\x0b\x61ttachments\x12\x1a\n\x08gateways\x18\x06 \x03(\tR\x08gateways\x12\x17\n\x04\x63hat\x18\x07 \x01(\tH\x00R\x04\x63hat\x88\x01\x01\x12\x12\n\x04sent\x18\x08 \x01(\x03R\x04sent\x12\x16\n\x06\x65\x64ited\x18\t \x01(\x03R\x06\x65\x64ited\x12!\n\x0cunder_review\x18\n \x01(\x08R\x0bunderReview\x12\x18\n\x07readers\x18\x0b \x03(\tR\x07readers\x12)\n\x04meta\x18\x0c \x03(\x0b\x32\x15.cc.Message.MetaEntryR\x04meta\x12\x1c\n\tmentioned\x18\r \x03(\tR\tmentioned\x1aO\n\tMetaEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12,\n\x05value\x18\x02 \x01(\x0b\x32\x16.google.protobuf.ValueR\x05value:\x02\x38\x01\x42\x07\n\x05_chat\"3\n\x08Messages\x12\'\n\x08messages\x18\x01 \x03(\x0b\x32\x0b.cc.MessageR\x08messages\"\x94\x02\n\x04User\x12\x12\n\x04uuid\x18\x01 \x01(\tR\x04uuid\x12\x14\n\x05title\x18\x02 \x01(\tR\x05title\x12+\n\x04\x64\x61ta\x18\x03 \x01(\x0b\x32\x17.google.protobuf.StructR\x04\x64\x61ta\x12\x1a\n\tcc_is_bot\x18\x04 \x01(\x08R\x07\x63\x63IsBot\x12\x1f\n\x0b\x63\x63_username\x18\x05 \x01(\tR\nccUsername\x12\x39\n\x0b\x63\x63_commands\x18\x06 \x03(\x0b\x32\x18.cc.User.CcCommandsEntryR\nccCommands\x1a=\n\x0f\x43\x63\x43ommandsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\"\xb3\x01\n\nDepartment\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05title\x18\x02 \x01(\tR\x05title\x12 \n\x0b\x64\x65scription\x18\x03 \x01(\tR\x0b\x64\x65scription\x12\x16\n\x06\x61\x64mins\x18\x04 \x03(\tR\x06\x61\x64mins\x12\x16\n\x06public\x18\x05 \x01(\x08R\x06public\x12\x1e\n\x08whmcs_id\x18\x06 \x01(\tH\x00R\x07whmcsId\x88\x01\x01\x42\x0b\n\t_whmcs_id\"0\n\x06Option\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\x02R\x05value\"D\n\x06Metric\x12\x14\n\x05title\x18\x01 \x01(\tR\x05title\x12$\n\x07options\x18\x02 \x03(\x0b\x32\n.cc.OptionR\x07options\"\x81\x03\n\x08\x44\x65\x66\x61ults\x12\x1a\n\x08gateways\x18\x01 \x03(\tR\x08gateways\x12\x16\n\x06\x61\x64mins\x18\x02 \x03(\tR\x06\x61\x64mins\x12\x30\n\x0b\x64\x65partments\x18\x03 \x03(\x0b\x32\x0e.cc.DepartmentR\x0b\x64\x65partments\x12\x33\n\x07metrics\x18\x04 \x03(\x0b\x32\x19.cc.Defaults.MetricsEntryR\x07metrics\x12\x39\n\ttemplates\x18\x05 \x03(\x0b\x32\x1b.cc.Defaults.TemplatesEntryR\ttemplates\x12\x19\n\x03\x62ot\x18\x06 \x01(\x0b\x32\x07.cc.BotR\x03\x62ot\x1a\x46\n\x0cMetricsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12 \n\x05value\x18\x02 \x01(\x0b\x32\n.cc.MetricR\x05value:\x02\x38\x01\x1a<\n\x0eTemplatesEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\"\x90\x02\n\x03\x42ot\x12\x16\n\x06prompt\x18\x01 \x01(\tR\x06prompt\x12:\n\x1a\x65nable_review_in_new_chats\x18\x02 \x01(\x08R\x16\x65nableReviewInNewChats\x12\x34\n\x17\x65nable_bot_in_new_chats\x18\x03 \x01(\x08R\x13\x65nableBotInNewChats\x12>\n\rcustom_values\x18\x04 \x03(\x0b\x32\x19.cc.Bot.CustomValuesEntryR\x0c\x63ustomValues\x1a?\n\x11\x43ustomValuesEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\"\'\n\x05Users\x12\x1e\n\x05users\x18\x01 \x03(\x0b\x32\x08.cc.UserR\x05users\"?\n\x14\x46\x65tchDefaultsRequest\x12\'\n\x0f\x66\x65tch_templates\x18\x01 \x01(\x08R\x0e\x66\x65tchTemplates\")\n\x11GetMembersRequest\x12\x14\n\x05uuids\x18\x01 \x03(\tR\x05uuids\"s\n\x05\x45vent\x12!\n\x04type\x18\x01 \x01(\x0e\x32\r.cc.EventTypeR\x04type\x12\x1e\n\x04\x63hat\x18\x02 \x01(\x0b\x32\x08.cc.ChatH\x00R\x04\x63hat\x12\x1f\n\x03msg\x18\x03 \x01(\x0b\x32\x0b.cc.MessageH\x00R\x03msgB\x06\n\x04item\"\x89\x01\n\rStreamRequest\x12;\n\x08\x63ommands\x18\x01 \x03(\x0b\x32\x1f.cc.StreamRequest.CommandsEntryR\x08\x63ommands\x1a;\n\rCommandsEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01*4\n\x04Role\x12\x0c\n\x08NOACCESS\x10\x00\x12\x08\n\x04USER\x10\x01\x12\t\n\x05OWNER\x10\x02\x12\t\n\x05\x41\x44MIN\x10\x03*\x8a\x01\n\x06Status\x12\x07\n\x03NEW\x10\x00\x12\x08\n\x04OPEN\x10\x01\x12\x0b\n\x07RESOLVE\x10\x02\x12\t\n\x05\x43LOSE\x10\x03\x12\x0c\n\x08\x41NSWERED\x10\x04\x12\x12\n\x0e\x43USTOMER_REPLY\x10\x05\x12\x15\n\x11WAITING_FOR_REPLY\x10\x06\x12\x0b\n\x07ON_HOLD\x10\x07\x12\x0f\n\x0bIN_PROGRESS\x10\x08*0\n\x04Kind\x12\x0b\n\x07\x44\x45\x46\x41ULT\x10\x00\x12\x0e\n\nADMIN_ONLY\x10\x01\x12\x0b\n\x07\x46OR_BOT\x10\x02*\xde\x01\n\tEventType\x12\x08\n\x04PING\x10\x00\x12\x10\n\x0c\x43HAT_CREATED\x10\x01\x12\x10\n\x0c\x43HAT_UPDATED\x10\x02\x12\x10\n\x0c\x43HAT_DELETED\x10\x03\x12\r\n\tCHAT_READ\x10\x04\x12\x10\n\x0cMESSAGE_SENT\x10\x05\x12\x13\n\x0fMESSAGE_UPDATED\x10\x06\x12\x13\n\x0fMESSAGE_DELETED\x10\x07\x12\x1b\n\x17\x43HAT_DEPARTMENT_CHANGED\x10\x08\x12\x17\n\x13\x43HAT_STATUS_CHANGED\x10\t\x12\x10\n\x0c\x43HATS_MERGED\x10\n2\xf7\x03\n\x08\x43hatsAPI\x12\x1c\n\x06\x43reate\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12\x1c\n\x06Update\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12\x19\n\x03Get\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12\x33\n\x04List\x12\x14.cc.ListChatsRequest\x1a\x15.cc.ListChatsResponse\x12\x36\n\x05\x43ount\x12\x15.cc.CountChatsRequest\x1a\x16.cc.CountChatsResponse\x12\x1c\n\x06\x44\x65lete\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12/\n\x0bSetBotState\x12\x16.cc.SetBotStateRequest\x1a\x08.cc.Chat\x12!\n\x0bGetBotState\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12&\n\x10\x43hangeDepartment\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12#\n\rChangeGateway\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12\"\n\x0c\x43hangeStatus\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12!\n\nMergeChats\x12\t.cc.Merge\x1a\x08.cc.Chat\x12!\n\tSyncChats\x12\t.cc.Empty\x1a\t.cc.Empty2\x96\x01\n\x0bMessagesAPI\x12\x1d\n\x03Get\x12\x08.cc.Chat\x1a\x0c.cc.Messages\x12 \n\x04Send\x12\x0b.cc.Message\x1a\x0b.cc.Message\x12\"\n\x06Update\x12\x0b.cc.Message\x1a\x0b.cc.Message\x12\"\n\x06\x44\x65lete\x12\x0b.cc.Message\x1a\x0b.cc.Message2\xfe\x01\n\x08UsersAPI\x12\x19\n\x02Me\x12\t.cc.Empty\x1a\x08.cc.User\x12\x37\n\rFetchDefaults\x12\x18.cc.FetchDefaultsRequest\x1a\x0c.cc.Defaults\x12$\n\tGetConfig\x12\t.cc.Empty\x1a\x0c.cc.Defaults\x12\'\n\tSetConfig\x12\x0c.cc.Defaults\x1a\x0c.cc.Defaults\x12\x1f\n\x07Resolve\x12\t.cc.Users\x1a\t.cc.Users\x12.\n\nGetMembers\x12\x15.cc.GetMembersRequest\x1a\t.cc.Users29\n\rStreamService\x12(\n\x06Stream\x12\x11.cc.StreamRequest\x1a\t.cc.Event0\x01\x42%Z#github.com/slntopp/core-chatting/ccb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cc.cc_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z#github.com/slntopp/core-chatting/cc'
  _globals['_CHATMETA_DATAENTRY']._loaded_options = None
  _globals['_CHATMETA_DATAENTRY']._serialized_options = b'8\001'
  _globals['_CHAT_BOTSTATEENTRY']._loaded_options = None
  _globals['_CHAT_BOTSTATEENTRY']._serialized_options = b'8\001'
  _globals['_LISTCHATSREQUEST_FILTERSENTRY']._loaded_options = None
  _globals['_LISTCHATSREQUEST_FILTERSENTRY']._serialized_options = b'8\001'
  _globals['_SETBOTSTATEREQUEST_STATEENTRY']._loaded_options = None
  _globals['_SETBOTSTATEREQUEST_STATEENTRY']._serialized_options = b'8\001'
  _globals['_COUNTCHATSREQUEST_FILTERSENTRY']._loaded_options = None
  _globals['_COUNTCHATSREQUEST_FILTERSENTRY']._serialized_options = b'8\001'
  _globals['_COUNTCHATSRESPONSE_STATUSESENTRY']._loaded_options = None
  _globals['_COUNTCHATSRESPONSE_STATUSESENTRY']._serialized_options = b'8\001'
  _globals['_MESSAGE_METAENTRY']._loaded_options = None
  _globals['_MESSAGE_METAENTRY']._serialized_options = b'8\001'
  _globals['_USER_CCCOMMANDSENTRY']._loaded_options = None
  _globals['_USER_CCCOMMANDSENTRY']._serialized_options = b'8\001'
  _globals['_DEFAULTS_METRICSENTRY']._loaded_options = None
  _globals['_DEFAULTS_METRICSENTRY']._serialized_options = b'8\001'
  _globals['_DEFAULTS_TEMPLATESENTRY']._loaded_options = None
  _globals['_DEFAULTS_TEMPLATESENTRY']._serialized_options = b'8\001'
  _globals['_BOT_CUSTOMVALUESENTRY']._loaded_options = None
  _globals['_BOT_CUSTOMVALUESENTRY']._serialized_options = b'8\001'
  _globals['_STREAMREQUEST_COMMANDSENTRY']._loaded_options = None
  _globals['_STREAMREQUEST_COMMANDSENTRY']._serialized_options = b'8\001'
  _globals['_ROLE']._serialized_start=4177
  _globals['_ROLE']._serialized_end=4229
  _globals['_STATUS']._serialized_start=4232
  _globals['_STATUS']._serialized_end=4370
  _globals['_KIND']._serialized_start=4372
  _globals['_KIND']._serialized_end=4420
  _globals['_EVENTTYPE']._serialized_start=4423
  _globals['_EVENTTYPE']._serialized_end=4645
  _globals['_EMPTY']._serialized_start=49
  _globals['_EMPTY']._serialized_end=56
  _globals['_CHATMETA']._serialized_start=59
  _globals['_CHATMETA']._serialized_end=316
  _globals['_CHATMETA_DATAENTRY']._serialized_start=237
  _globals['_CHATMETA_DATAENTRY']._serialized_end=316
  _globals['_CHAT']._serialized_start=319
  _globals['_CHAT']._serialized_end=843
  _globals['_CHAT_BOTSTATEENTRY']._serialized_start=725
  _globals['_CHAT_BOTSTATEENTRY']._serialized_end=808
  _globals['_CHATS']._serialized_start=845
  _globals['_CHATS']._serialized_end=884
  _globals['_MERGE']._serialized_start=886
  _globals['_MERGE']._serialized_end=915
  _globals['_LISTCHATSREQUEST']._serialized_start=918
  _globals['_LISTCHATSREQUEST']._serialized_end=1223
  _globals['_LISTCHATSREQUEST_FILTERSENTRY']._serialized_start=1103
  _globals['_LISTCHATSREQUEST_FILTERSENTRY']._serialized_end=1185
  _globals['_SETBOTSTATEREQUEST']._serialized_start=1226
  _globals['_SETBOTSTATEREQUEST']._serialized_end=1554
  _globals['_SETBOTSTATEREQUEST_STATEENTRY']._serialized_start=1431
  _globals['_SETBOTSTATEREQUEST_STATEENTRY']._serialized_end=1511
  _globals['_LISTCHATSRESPONSE']._serialized_start=1556
  _globals['_LISTCHATSRESPONSE']._serialized_end=1627
  _globals['_COUNTCHATSREQUEST']._serialized_start=1630
  _globals['_COUNTCHATSREQUEST']._serialized_end=1795
  _globals['_COUNTCHATSREQUEST_FILTERSENTRY']._serialized_start=1103
  _globals['_COUNTCHATSREQUEST_FILTERSENTRY']._serialized_end=1185
  _globals['_COUNTCHATSRESPONSE']._serialized_start=1798
  _globals['_COUNTCHATSRESPONSE']._serialized_end=1945
  _globals['_COUNTCHATSRESPONSE_STATUSESENTRY']._serialized_start=1886
  _globals['_COUNTCHATSRESPONSE_STATUSESENTRY']._serialized_end=1945
  _globals['_ATTACHMENT']._serialized_start=1947
  _globals['_ATTACHMENT']._serialized_end=2005
  _globals['_MESSAGE']._serialized_start=2008
  _globals['_MESSAGE']._serialized_end=2472
  _globals['_MESSAGE_METAENTRY']._serialized_start=2384
  _globals['_MESSAGE_METAENTRY']._serialized_end=2463
  _globals['_MESSAGES']._serialized_start=2474
  _globals['_MESSAGES']._serialized_end=2525
  _globals['_USER']._serialized_start=2528
  _globals['_USER']._serialized_end=2804
  _globals['_USER_CCCOMMANDSENTRY']._serialized_start=2743
  _globals['_USER_CCCOMMANDSENTRY']._serialized_end=2804
  _globals['_DEPARTMENT']._serialized_start=2807
  _globals['_DEPARTMENT']._serialized_end=2986
  _globals['_OPTION']._serialized_start=2988
  _globals['_OPTION']._serialized_end=3036
  _globals['_METRIC']._serialized_start=3038
  _globals['_METRIC']._serialized_end=3106
  _globals['_DEFAULTS']._serialized_start=3109
  _globals['_DEFAULTS']._serialized_end=3494
  _globals['_DEFAULTS_METRICSENTRY']._serialized_start=3362
  _globals['_DEFAULTS_METRICSENTRY']._serialized_end=3432
  _globals['_DEFAULTS_TEMPLATESENTRY']._serialized_start=3434
  _globals['_DEFAULTS_TEMPLATESENTRY']._serialized_end=3494
  _globals['_BOT']._serialized_start=3497
  _globals['_BOT']._serialized_end=3769
  _globals['_BOT_CUSTOMVALUESENTRY']._serialized_start=3706
  _globals['_BOT_CUSTOMVALUESENTRY']._serialized_end=3769
  _globals['_USERS']._serialized_start=3771
  _globals['_USERS']._serialized_end=3810
  _globals['_FETCHDEFAULTSREQUEST']._serialized_start=3812
  _globals['_FETCHDEFAULTSREQUEST']._serialized_end=3875
  _globals['_GETMEMBERSREQUEST']._serialized_start=3877
  _globals['_GETMEMBERSREQUEST']._serialized_end=3918
  _globals['_EVENT']._serialized_start=3920
  _globals['_EVENT']._serialized_end=4035
  _globals['_STREAMREQUEST']._serialized_start=4038
  _globals['_STREAMREQUEST']._serialized_end=4175
  _globals['_STREAMREQUEST_COMMANDSENTRY']._serialized_start=4116
  _globals['_STREAMREQUEST_COMMANDSENTRY']._serialized_end=4175
  _globals['_CHATSAPI']._serialized_start=4648
  _globals['_CHATSAPI']._serialized_end=5151
  _globals['_MESSAGESAPI']._serialized_start=5154
  _globals['_MESSAGESAPI']._serialized_end=5304
  _globals['_USERSAPI']._serialized_start=5307
  _globals['_USERSAPI']._serialized_end=5561
  _globals['_STREAMSERVICE']._serialized_start=5563
  _globals['_STREAMSERVICE']._serialized_end=5620
# @@protoc_insertion_point(module_scope)
