<script lang="ts">
    import { login } from "../lib/auth";

    let isLoggingIn = false;
    let error = "";

    async function handleLogin() {
        isLoggingIn = true;
        error = "";
        try {
            await login();
        } catch (e) {
            error = e instanceof Error ? e.message : "Failed to initiate login";
            isLoggingIn = false;
        }
    }
</script>

<div
    class="min-h-screen flex items-center justify-center bg-gradient-to-br from-red-900 via-gray-900 to-gray-800"
>
    <div class="w-full max-w-md p-8">
        <!-- Logo/Branding -->
        <div class="text-center mb-8">
            <div
                class="inline-flex items-center justify-center w-20 h-20 bg-red-600 rounded-full mb-4"
            >
                <svg
                    class="w-10 h-10 text-white"
                    fill="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"
                    />
                </svg>
            </div>
            <h1 class="text-4xl font-bold text-white mb-2">CMU Insta</h1>
            <p class="text-gray-400 text-lg">Share your CMU moments</p>
        </div>

        <!-- Login Card -->
        <div class="bg-white rounded-2xl shadow-2xl p-8">
            <div class="text-center mb-6">
                <h2 class="text-2xl font-semibold text-gray-800 mb-2">
                    Welcome
                </h2>
                <p class="text-gray-600">
                    Sign in with your Andrew ID to continue
                </p>
            </div>

            {#if error}
                <div
                    class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg"
                >
                    <p class="text-red-700 text-sm text-center">{error}</p>
                </div>
            {/if}

            <button
                on:click={handleLogin}
                disabled={isLoggingIn}
                class="w-full py-4 px-6 bg-red-600 hover:bg-red-700 disabled:bg-red-400
               text-white font-semibold rounded-xl transition-all duration-200
               flex items-center justify-center gap-3 shadow-lg hover:shadow-xl
               disabled:cursor-not-allowed"
            >
                {#if isLoggingIn}
                    <svg
                        class="animate-spin h-5 w-5 text-white"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                    >
                        <circle
                            class="opacity-25"
                            cx="12"
                            cy="12"
                            r="10"
                            stroke="currentColor"
                            stroke-width="4"
                        ></circle>
                        <path
                            class="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        ></path>
                    </svg>
                    <span>Redirecting to CMU SSO...</span>
                {:else}
                    <svg
                        class="w-5 h-5"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"
                        />
                    </svg>
                    <span>Sign in with Andrew ID</span>
                {/if}
            </button>

            <div class="mt-6 text-center">
                <p class="text-xs text-gray-500">
                    By signing in, you agree to CMU's policies and guidelines.
                </p>
            </div>
        </div>

        <!-- Footer -->
        <div class="text-center mt-8">
            <p class="text-gray-500 text-sm">
                Powered by <span class="text-red-400 font-medium"
                    >ScottyLabs</span
                > SSO
            </p>
            <p class="text-gray-600 text-xs mt-2">Carnegie Mellon University</p>
        </div>
    </div>
</div>
