import Vue from 'vue'
import Router from 'vue-router'
import auth from './auth'
import Home from './views/Home.vue'
import Login from './views/Login.vue'
import Register from './views/Register.vue'
import NewPost from './views/NewPost.vue'
import Profile from './views/Profile.vue'
import CommentPage from './views/CommentPage.vue'
import Subnoddit from './views/Subnoddit.vue'
import NewSubnoddit from './views/NewSubnoddit.vue'
import AllSubnoddits from './views/AllSubnoddits.vue'
import SubnodditSearch from './views/SubnodditSearch.vue'

Vue.use(Router)

/**
 * The Vue Router is used to "direct" the browser to render a specific view component
 * inside of App.vue depending on the URL.
 *
 * It also is used to detect whether or not a route requires the user to have first authenticated.
 * If the user has not yet authenticated (and needs to) they are redirected to /login
 * If they have (or don't need to) they're allowed to go about their way.
 */

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
      meta: {
        requiresAuth: false
      }
    },
    {
      path: "/login",
      name: "login",
      component: Login,
      meta: {
        requiresAuth: false
      }
    },
    {
      path: "/register",
      name: "register",
      component: Register,
      meta: {
        requiresAuth: false
      }
    },
    {
      path: "/newpost",
      name: "newpost",
      component: NewPost,
      meta: {
        requiresAuth: true
      }
    },
    {
      path: "/n/:subnoddit",
      name: "subnoddit",
      component: Subnoddit,
      meta: {
        requiresAuth: false
      }
    },
    {
      path: "/n/:subnoddit/:id",
      name: "subnoddit-post",
      component: CommentPage,
      meta: {
        requiresAuth: false
      }
    },
    {
      path: "/profile",
      name: "profile",
      component: Profile,
      meta: {
        requiresAuth: true
      }
    },
    {
      path: "/newSubnoddit",
      name: "newSubnoddit",
      component: NewSubnoddit,
      meta: {
        requiresAuth: true
      }
    },
    {
      path: "/allsubnoddits",
      name: "allsubnoddits",
      component: AllSubnoddits,
      meta: {
        requiresAuth: false
      }
    },
    {
      path: "/search/:searchTerm",
      name: "search",
      component: SubnodditSearch,
      meta: {
        requiresAuth: false
      }
    },
    {
      path: "/search",
      name: "search",
      component: SubnodditSearch,
      meta: {
        requiresAuth: false
      }
    }
  ]
})

router.beforeEach((to, from, next) => {
  // Determine if the route requires Authentication
  const requiresAuth = to.matched.some(x => x.meta.requiresAuth);
  const user = auth.getUser();

  // If it does and they are not logged in, send the user to "/login"
  if (requiresAuth && !user) {
    next("/login");
  } else {
    // Else let them go to their next destination
    next();
  }
});

export default router;
