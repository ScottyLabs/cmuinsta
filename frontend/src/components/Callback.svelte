<script lang="ts">
    import { onMount } from "svelte";
    import { handleCallback } from "../lib/auth";
    import { authStore, currentRoute } from "../lib/stores";

    let error: string | null = null;
    let loading = true;

    onMount(async () => {
        try {
            const result = await handleCallback();

            if (result) {
                // Successfully authenticated
                authStore.setUser(
                    {
                        andrewId: result.user.andrewId,
                        name:
                            result.user.name ||
                            result.user.givenName ||
                            result.user.andrewId,
                        email: result.user.email,
                        isAdmin: result.isAdmin,
                    },
                    localStorage.getItem("auth_access_token") || "",
                );

                // Redirect based on admin status
                if (result.isAdmin) {
                    currentRoute.set("admin");
                } else {
                    currentRoute.set("home");
                }
            } else {
                // No code in URL, redirect to login
                currentRoute.set("login");
            }
        } catch (e) {
            console.error("Callback error:", e);
            error = e instanceof Error ? e.message : "Authentication failed";
            authStore.setError(error);
        } finally {
            loading = false;
        }
    });
</script>

<main
    class="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-red-900 via-gray-900 to-black"
>
    <div
        class="w-full max-w-md rounded-xl bg-white/10 backdrop-blur-md p-8 shadow-2xl border border-white/20"
    >
        {#if loading}
            <div class="flex flex-col items-center space-y-4">
                <div
                    class="w-12 h-12 border-4 border-white/30 border-t-white rounded-full animate-spin"
                ></div>
                <h2 class="text-xl font-semibold text-white">
                    Authenticating...
                </h2>
                <p class="text-gray-300 text-sm">
                    Please wait while we verify your CMU credentials
                </p>
            </div>
        {:else if error}
            <div class="flex flex-col items-center space-y-4">
                <div
                    class="w-16 h-16 bg-red-500/20 rounded-full flex items-center justify-center"
                >
                    <svg
                        class="w-8 h-8 text-red-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M6 18L18 6M6 6l12 12"
                        ></path>
                    </svg>
                </div>
                <h2 class="text-xl font-semibold text-white">
                    Authentication Failed
                </h2>
                <p class="text-red-300 text-sm text-center">{error}</p>
                <button
                    on:click={() => currentRoute.set("login")}
                    class="mt-4 px-6 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
                >
                    Back to Login
                </button>
            </div>
        {/if}
    </div>
</main>
