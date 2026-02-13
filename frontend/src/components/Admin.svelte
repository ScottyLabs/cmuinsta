<script lang="ts">
    import { authStore, currentRoute } from "../lib/stores";
    import { logout } from "../lib/auth";

    let andrewId = "";
    let name = "";

    authStore.subscribe((state) => {
        if (state.user) {
            andrewId = state.user.andrewId;
            name = state.user.name;
        }
    });

    function handleLogout() {
        logout();
    }

    function goToHome() {
        currentRoute.set("home");
    }
</script>

<div class="min-h-screen bg-gray-900">
    <!-- Header -->
    <header class="bg-gray-800 shadow-lg border-b border-gray-700">
        <div
            class="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between"
        >
            <div class="flex items-center gap-3">
                <h1 class="text-2xl font-bold text-red-500">CMU Insta</h1>
                <span
                    class="px-2 py-1 bg-red-600 text-white text-xs font-semibold rounded-md uppercase"
                >
                    Admin
                </span>
            </div>
            <div class="flex items-center gap-4">
                <span class="text-gray-400 text-sm">
                    <span class="font-medium text-white">{name}</span>
                    <span class="text-gray-500">({andrewId})</span>
                </span>
                <button
                    on:click={goToHome}
                    class="px-4 py-2 text-sm text-gray-400 hover:text-white transition-colors"
                >
                    User View
                </button>
                <button
                    on:click={handleLogout}
                    class="px-4 py-2 text-sm text-gray-400 hover:text-red-400 transition-colors"
                >
                    Sign Out
                </button>
            </div>
        </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-6xl mx-auto px-4 py-12">
        <div
            class="bg-gray-800 rounded-2xl shadow-2xl p-12 border border-gray-700"
        >
            <div class="text-center">
                <!-- Icon -->
                <div
                    class="inline-flex items-center justify-center w-24 h-24 bg-gray-700 rounded-full mb-6"
                >
                    <svg
                        class="w-12 h-12 text-gray-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                        ></path>
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                        ></path>
                    </svg>
                </div>

                <h2 class="text-3xl font-bold text-white mb-4">
                    Admin Console
                </h2>
                <p class="text-gray-400 text-lg mb-8 max-w-md mx-auto">
                    This admin console is under construction. Check back later
                    for post moderation and management features.
                </p>

                <!-- Placeholder Stats -->
                <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
                    <div class="bg-gray-700/50 rounded-xl p-6">
                        <div class="text-4xl font-bold text-white mb-2">--</div>
                        <div class="text-gray-400 text-sm">Pending Posts</div>
                    </div>
                    <div class="bg-gray-700/50 rounded-xl p-6">
                        <div class="text-4xl font-bold text-white mb-2">--</div>
                        <div class="text-gray-400 text-sm">Approved Posts</div>
                    </div>
                    <div class="bg-gray-700/50 rounded-xl p-6">
                        <div class="text-4xl font-bold text-white mb-2">--</div>
                        <div class="text-gray-400 text-sm">Total Users</div>
                    </div>
                </div>

                <!-- Coming Soon Features -->
                <div class="mt-12 text-left">
                    <h3 class="text-lg font-semibold text-gray-300 mb-4">
                        Coming Soon
                    </h3>
                    <ul class="space-y-3">
                        <li class="flex items-center gap-3 text-gray-400">
                            <svg
                                class="w-5 h-5 text-gray-500"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                ></path>
                            </svg>
                            Review and approve/reject pending posts
                        </li>
                        <li class="flex items-center gap-3 text-gray-400">
                            <svg
                                class="w-5 h-5 text-gray-500"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                ></path>
                            </svg>
                            View all submitted posts with media previews
                        </li>
                        <li class="flex items-center gap-3 text-gray-400">
                            <svg
                                class="w-5 h-5 text-gray-500"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                ></path>
                            </svg>
                            Manage admin users and permissions
                        </li>
                        <li class="flex items-center gap-3 text-gray-400">
                            <svg
                                class="w-5 h-5 text-gray-500"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                ></path>
                            </svg>
                            Analytics and usage statistics
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </main>

    <!-- Footer -->
    <footer class="text-center py-6 text-gray-500 text-sm">
        CMU Insta Admin Console • Powered by ScottyLabs
    </footer>
</div>
