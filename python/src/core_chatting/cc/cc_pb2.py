# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: cc/cc.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import struct_pb2 as google_dot_protobuf_dot_struct__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0b\x63\x63/cc.proto\x12\x02\x63\x63\x1a\x1cgoogle/protobuf/struct.proto\"\x07\n\x05\x45mpty\"R\n\x08\x43hatMeta\x12\x16\n\x06unread\x18\x01 \x01(\x05R\x06unread\x12.\n\x0clast_message\x18\x02 \x01(\x0b\x32\x0b.cc.MessageR\x0blastMessage\"\x87\x02\n\x04\x43hat\x12\x12\n\x04uuid\x18\x01 \x01(\tR\x04uuid\x12\x14\n\x05users\x18\x02 \x03(\tR\x05users\x12\x16\n\x06\x61\x64mins\x18\x03 \x03(\tR\x06\x61\x64mins\x12\x14\n\x05owner\x18\x04 \x01(\tR\x05owner\x12\x1a\n\x08gateways\x18\x05 \x03(\tR\x08gateways\x12\x1c\n\x04role\x18\x06 \x01(\x0e\x32\x08.cc.RoleR\x04role\x12\x19\n\x05topic\x18\x07 \x01(\tH\x00R\x05topic\x88\x01\x01\x12%\n\x04meta\x18\x08 \x01(\x0b\x32\x0c.cc.ChatMetaH\x01R\x04meta\x88\x01\x01\x12\x18\n\x07\x63reated\x18\t \x01(\x03R\x07\x63reatedB\x08\n\x06_topicB\x07\n\x05_meta\"\'\n\x05\x43hats\x12\x1e\n\x05\x63hats\x18\x01 \x03(\x0b\x32\x08.cc.ChatR\x05\x63hats\":\n\nAttachment\x12\x12\n\x04type\x18\x01 \x01(\tR\x04type\x12\x18\n\x07\x63ontent\x18\x02 \x01(\x0cR\x07\x63ontent\"\xc6\x02\n\x07Message\x12\x12\n\x04uuid\x18\x01 \x01(\tR\x04uuid\x12\x1c\n\x04kind\x18\x02 \x01(\x0e\x32\x08.cc.KindR\x04kind\x12\x16\n\x06sender\x18\x03 \x01(\tR\x06sender\x12\x18\n\x07\x63ontent\x18\x04 \x01(\tR\x07\x63ontent\x12\x30\n\x0b\x61ttachments\x18\x05 \x03(\x0b\x32\x0e.cc.AttachmentR\x0b\x61ttachments\x12\x1a\n\x08gateways\x18\x06 \x03(\tR\x08gateways\x12\x17\n\x04\x63hat\x18\x07 \x01(\tH\x00R\x04\x63hat\x88\x01\x01\x12\x12\n\x04sent\x18\x08 \x01(\x03R\x04sent\x12\x16\n\x06\x65\x64ited\x18\t \x01(\x03R\x06\x65\x64ited\x12!\n\x0cunder_review\x18\n \x01(\x08R\x0bunderReview\x12\x18\n\x07readers\x18\x0b \x03(\tR\x07readersB\x07\n\x05_chat\"3\n\x08Messages\x12\'\n\x08messages\x18\x01 \x03(\x0b\x32\x0b.cc.MessageR\x08messages\"0\n\x04User\x12\x12\n\x04uuid\x18\x01 \x01(\tR\x04uuid\x12\x14\n\x05title\x18\x02 \x01(\tR\x05title\">\n\x08\x44\x65\x66\x61ults\x12\x1a\n\x08gateways\x18\x01 \x03(\tR\x08gateways\x12\x16\n\x06\x61\x64mins\x18\x02 \x03(\tR\x06\x61\x64mins\"\'\n\x05Users\x12\x1e\n\x05users\x18\x01 \x03(\x0b\x32\x08.cc.UserR\x05users\"s\n\x05\x45vent\x12!\n\x04type\x18\x01 \x01(\x0e\x32\r.cc.EventTypeR\x04type\x12\x1e\n\x04\x63hat\x18\x02 \x01(\x0b\x32\x08.cc.ChatH\x00R\x04\x63hat\x12\x1f\n\x03msg\x18\x03 \x01(\x0b\x32\x0b.cc.MessageH\x00R\x03msgB\x06\n\x04item*4\n\x04Role\x12\x0c\n\x08NOACCESS\x10\x00\x12\x08\n\x04USER\x10\x01\x12\t\n\x05OWNER\x10\x02\x12\t\n\x05\x41\x44MIN\x10\x03*#\n\x04Kind\x12\x0b\n\x07\x44\x45\x46\x41ULT\x10\x00\x12\x0e\n\nADMIN_ONLY\x10\x01*\x87\x01\n\tEventType\x12\x08\n\x04PING\x10\x00\x12\x10\n\x0c\x43HAT_CREATED\x10\x01\x12\x10\n\x0c\x43HAT_UPDATED\x10\x02\x12\x10\n\x0c\x43HAT_DELETED\x10\x03\x12\x10\n\x0cMESSAGE_SENT\x10\x04\x12\x13\n\x0fMESSAGE_UPDATED\x10\x05\x12\x13\n\x0fMESSAGE_DELETED\x10\x06\x32\x9d\x01\n\x08\x43hatsAPI\x12\x1c\n\x06\x43reate\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12\x1c\n\x06Update\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12\x19\n\x03Get\x12\x08.cc.Chat\x1a\x08.cc.Chat\x12\x1c\n\x04List\x12\t.cc.Empty\x1a\t.cc.Chats\x12\x1c\n\x06\x44\x65lete\x12\x08.cc.Chat\x1a\x08.cc.Chat2\x96\x01\n\x0bMessagesAPI\x12\x1d\n\x03Get\x12\x08.cc.Chat\x1a\x0c.cc.Messages\x12 \n\x04Send\x12\x0b.cc.Message\x1a\x0b.cc.Message\x12\"\n\x06Update\x12\x0b.cc.Message\x1a\x0b.cc.Message\x12\"\n\x06\x44\x65lete\x12\x0b.cc.Message\x1a\x0b.cc.Message2p\n\x08UsersAPI\x12\x19\n\x02Me\x12\t.cc.Empty\x1a\x08.cc.User\x12(\n\rFetchDefaults\x12\t.cc.Empty\x1a\x0c.cc.Defaults\x12\x1f\n\x07Resolve\x12\t.cc.Users\x1a\t.cc.Users21\n\rStreamService\x12 \n\x06Stream\x12\t.cc.Empty\x1a\t.cc.Event0\x01\x42%Z#github.com/slntopp/core-chatting/ccb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'cc.cc_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z#github.com/slntopp/core-chatting/cc'
  _globals['_ROLE']._serialized_start=1163
  _globals['_ROLE']._serialized_end=1215
  _globals['_KIND']._serialized_start=1217
  _globals['_KIND']._serialized_end=1252
  _globals['_EVENTTYPE']._serialized_start=1255
  _globals['_EVENTTYPE']._serialized_end=1390
  _globals['_EMPTY']._serialized_start=49
  _globals['_EMPTY']._serialized_end=56
  _globals['_CHATMETA']._serialized_start=58
  _globals['_CHATMETA']._serialized_end=140
  _globals['_CHAT']._serialized_start=143
  _globals['_CHAT']._serialized_end=406
  _globals['_CHATS']._serialized_start=408
  _globals['_CHATS']._serialized_end=447
  _globals['_ATTACHMENT']._serialized_start=449
  _globals['_ATTACHMENT']._serialized_end=507
  _globals['_MESSAGE']._serialized_start=510
  _globals['_MESSAGE']._serialized_end=836
  _globals['_MESSAGES']._serialized_start=838
  _globals['_MESSAGES']._serialized_end=889
  _globals['_USER']._serialized_start=891
  _globals['_USER']._serialized_end=939
  _globals['_DEFAULTS']._serialized_start=941
  _globals['_DEFAULTS']._serialized_end=1003
  _globals['_USERS']._serialized_start=1005
  _globals['_USERS']._serialized_end=1044
  _globals['_EVENT']._serialized_start=1046
  _globals['_EVENT']._serialized_end=1161
  _globals['_CHATSAPI']._serialized_start=1393
  _globals['_CHATSAPI']._serialized_end=1550
  _globals['_MESSAGESAPI']._serialized_start=1553
  _globals['_MESSAGESAPI']._serialized_end=1703
  _globals['_USERSAPI']._serialized_start=1705
  _globals['_USERSAPI']._serialized_end=1817
  _globals['_STREAMSERVICE']._serialized_start=1819
  _globals['_STREAMSERVICE']._serialized_end=1868
# @@protoc_insertion_point(module_scope)
