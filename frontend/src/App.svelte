<script lang="ts">
    import { onMount } from "svelte";
    import Login from "./components/Login.svelte";
    import Callback from "./components/Callback.svelte";
    import Home from "./components/Home.svelte";
    import Admin from "./components/Admin.svelte";
    import { authStore, currentRoute, type Route } from "./lib/stores";
    import { getAuthState } from "./lib/auth";

    let isInitialized = false;

    onMount(() => {
        // Check current URL path for routing
        const path = window.location.pathname;

        if (
            path === "/oauth2/callback" ||
            path.startsWith("/oauth2/callback")
        ) {
            currentRoute.set("callback");
        } else {
            // Check if user is already authenticated
            const authState = getAuthState();

            if (authState.isAuthenticated && authState.user) {
                // Restore auth state from storage
                authStore.setUser(
                    {
                        andrewId: authState.user.andrewId,
                        name:
                            authState.user.name ||
                            authState.user.givenName ||
                            authState.user.andrewId,
                        email: authState.user.email,
                        isAdmin: authState.isAdmin,
                    },
                    authState.accessToken || "",
                );

                // Route based on admin status
                if (authState.isAdmin) {
                    currentRoute.set("admin");
                } else {
                    currentRoute.set("home");
                }
            } else {
                authStore.setLoading(false);
                currentRoute.set("login");
            }
        }

        isInitialized = true;

        // Handle browser back/forward
        window.addEventListener("popstate", handlePopState);

        return () => {
            window.removeEventListener("popstate", handlePopState);
        };
    });

    function handlePopState() {
        const path = window.location.pathname;
        if (path === "/oauth2/callback") {
            currentRoute.set("callback");
        } else if (path === "/admin") {
            currentRoute.set("admin");
        } else if (path === "/home" || path === "/") {
            const state = getAuthState();
            if (state.isAuthenticated) {
                currentRoute.set(state.isAdmin ? "admin" : "home");
            } else {
                currentRoute.set("login");
            }
        }
    }

    // Update URL when route changes
    $: if (isInitialized && $currentRoute) {
        const pathMap: Record<Route, string> = {
            login: "/",
            callback: "/oauth2/callback",
            home: "/home",
            admin: "/admin",
        };
        const newPath = pathMap[$currentRoute];
        if (window.location.pathname !== newPath) {
            window.history.pushState({}, "", newPath);
        }
    }
</script>

{#if !isInitialized}
    <div
        class="min-h-screen flex items-center justify-center bg-gradient-to-br from-red-900 via-gray-900 to-gray-800"
    >
        <div class="flex flex-col items-center space-y-4">
            <div
                class="w-12 h-12 border-4 border-white/30 border-t-white rounded-full animate-spin"
            ></div>
            <p class="text-white text-lg">Loading...</p>
        </div>
    </div>
{:else if $currentRoute === "login"}
    <Login />
{:else if $currentRoute === "callback"}
    <Callback />
{:else if $currentRoute === "home"}
    <Home />
{:else if $currentRoute === "admin"}
    <Admin />
{/if}
