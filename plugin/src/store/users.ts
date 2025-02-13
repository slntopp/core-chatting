import { defineStore, storeToRefs } from 'pinia'
import { useCcStore } from './chatting';
import { ref, watch } from 'vue';
import { GetMembersRequest, User, Users } from '../connect/cc/cc_pb';
import { useDefaultsStore } from './defaults';

export const useUsersStore = defineStore('users', () => {
    const ccStore = useCcStore()
    const { chats, currentChat } = storeToRefs(ccStore)
    const defaultsStore = useDefaultsStore()
    const { admins } = storeToRefs(defaultsStore)

    const users = ref(new Map<string, User>());
    const usersForFetch = ref(new Map<string, boolean>());
    const isUsersLoading = ref(false);


    async function fetchUsers(uuids: string[]) {
        const resultUuids = []
        for (const uuid of uuids) {
            if (!usersForFetch.value.get(uuid)) {
                usersForFetch.value.set(uuid, true)
                resultUuids.push(uuid)
            }
        }

        if (resultUuids.length === 0) {
            return
        }

        isUsersLoading.value = true;

        try {
            const data = (await ccStore.users_c.getMembers(new GetMembersRequest({ uuids: resultUuids }))).toJson() as any as Users;

            data.users.forEach(user => {
                users.value.set(user.uuid, user);
            })
        } catch (e) {
            for (const uuid of uuids) {
                users.value.delete(uuid,)
                usersForFetch.value.delete(uuid,)
            }
        } finally {
            isUsersLoading.value = false;
        }
    }

    watch(chats, (value) => {
        fetchUsers(([...value.values()]).map((chat) => [...chat.admins, ...chat.users]).flat());
    });

    watch(admins, (value) => {
        fetchUsers(value);
    });

    watch(currentChat, (value) => {
        if (value) {
            fetchUsers([...value.admins, ...value.users]);
        }
    });

    return {
        users,
        isUsersLoading,
        fetchUsers
    }
})