import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/",
    name: "Entrypoint",
    component: () => import("../views/Entrypoint.vue"),
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("../views/Dashboard.vue"),
    children: [
      {
        path: '',
        name: 'Chats',
        component: () => import("../views/Dashboard/Chats.vue"),
        children: [
          {
            path: '',
            name: 'Empty Chat',
            component: () => import("../views/Dashboard/Chats/Empty.vue"),
          },
          {
            path: 'start',
            name: 'Start Chat',
            component: () => import("../views/Dashboard/Chats/Start.vue"),
          }
        ]
      }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;