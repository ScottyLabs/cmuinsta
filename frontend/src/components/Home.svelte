<script lang="ts">
    import { onMount } from "svelte";
    import { authStore, currentRoute } from "../lib/stores";
    import { logout, getAccessToken } from "../lib/auth";
    import type { UploadFile } from "../lib/types";

    // Form state
    let name = "";
    let andrewId = "";
    let instagramUsername = "";
    let caption = "";
    let uploadFiles: UploadFile[] = [];
    let isSubmitting = false;
    let submitError = "";
    let submitSuccess = false;

    // Instagram validation state
    let isValidatingInstagram = false;
    let instagramValid: boolean | null = null;
    let instagramError = "";

    // File input ref
    let fileInput: HTMLInputElement;

    // Character count
    $: charCount = caption.length;
    $: isOverLimit = charCount > 2200;

    // Clean Instagram username (remove @ if present)
    $: cleanInstagramUsername = instagramUsername.replace(/^@/, "").trim();

    // Debounced Instagram validation
    let instagramValidationTimeout: ReturnType<typeof setTimeout>;

    function handleInstagramInput() {
        // Reset validation state
        instagramValid = null;
        instagramError = "";

        // Clear previous timeout
        if (instagramValidationTimeout) {
            clearTimeout(instagramValidationTimeout);
        }

        // Don't validate empty input
        if (!cleanInstagramUsername) {
            return;
        }

        // Validate format first (alphanumeric, underscores, periods, 1-30 chars)
        const usernameRegex = /^[a-zA-Z0-9._]{1,30}$/;
        if (!usernameRegex.test(cleanInstagramUsername)) {
            instagramError = "Invalid username format";
            instagramValid = false;
            return;
        }

        // Debounce the API call
        instagramValidationTimeout = setTimeout(() => {
            validateInstagramUsername(cleanInstagramUsername);
        }, 500);
    }

    async function validateInstagramUsername(username: string) {
        isValidatingInstagram = true;
        instagramError = "";

        try {
            // Client-side validation only - we can't directly verify Instagram accounts
            // due to CORS restrictions, so we validate the format and show a preview
            // for the user to manually verify

            // Simulate a brief delay for UX
            await new Promise((resolve) => setTimeout(resolve, 300));

            // Format is already validated in handleInstagramInput
            // Mark as valid (format-wise) and show preview for user verification
            instagramValid = true;
            instagramError = "";
        } catch (e) {
            instagramValid = null;
            instagramError = "Could not verify username format";
        } finally {
            isValidatingInstagram = false;
        }
    }

    // Generate Instagram profile URL
    function getInstagramProfileUrl(username: string): string {
        return `https://www.instagram.com/${username}/`;
    }

    // Handle opening Instagram profile in new tab
    function openInstagramProfile() {
        if (cleanInstagramUsername) {
            window.open(
                getInstagramProfileUrl(cleanInstagramUsername),
                "_blank",
                "noopener,noreferrer",
            );
        }
    }

    onMount(() => {
        // Auto-fill from auth store
        const unsubscribe = authStore.subscribe((state) => {
            if (state.user) {
                andrewId = state.user.andrewId;
                name = state.user.name || "";
            }
        });

        return unsubscribe;
    });

    function handleFileSelect(event: Event) {
        const input = event.target as HTMLInputElement;
        if (!input.files) return;
        addFiles(Array.from(input.files));
        input.value = ""; // Reset input
    }

    function addFiles(files: File[]) {
        const remainingSlots = 10 - uploadFiles.length;
        const filesToAdd = files.slice(0, remainingSlots);

        filesToAdd.forEach((file) => {
            // Validate file type
            if (
                !file.type.startsWith("image/") &&
                !file.type.startsWith("video/")
            ) {
                return;
            }

            const id = crypto.randomUUID();
            const preview = URL.createObjectURL(file);
            uploadFiles = [
                ...uploadFiles,
                {
                    id,
                    file,
                    preview,
                    order: uploadFiles.length,
                },
            ];
        });
    }

    function removeFile(id: string) {
        const file = uploadFiles.find((f) => f.id === id);
        if (file) {
            URL.revokeObjectURL(file.preview);
        }
        uploadFiles = uploadFiles
            .filter((f) => f.id !== id)
            .map((f, idx) => ({ ...f, order: idx }));
    }

    // Drag and drop handlers
    let draggedItem: UploadFile | null = null;
    let dropTarget: HTMLDivElement;
    let isDraggingOver = false;

    function handleDragStart(event: DragEvent, file: UploadFile) {
        draggedItem = file;
        if (event.dataTransfer) {
            event.dataTransfer.effectAllowed = "move";
        }
    }

    function handleDragOver(event: DragEvent) {
        event.preventDefault();
        if (event.dataTransfer) {
            event.dataTransfer.dropEffect = "move";
        }
    }

    function handleDragEnter(event: DragEvent, targetFile: UploadFile) {
        event.preventDefault();
        if (!draggedItem || draggedItem.id === targetFile.id) return;

        // Reorder files
        const draggedIndex = uploadFiles.findIndex(
            (f) => f.id === draggedItem!.id,
        );
        const targetIndex = uploadFiles.findIndex(
            (f) => f.id === targetFile.id,
        );

        if (draggedIndex !== -1 && targetIndex !== -1) {
            const newFiles = [...uploadFiles];
            newFiles.splice(draggedIndex, 1);
            newFiles.splice(targetIndex, 0, draggedItem);
            uploadFiles = newFiles.map((f, idx) => ({ ...f, order: idx }));
        }
    }

    function handleDragEnd() {
        draggedItem = null;
    }

    // Drop zone handlers for new files
    function handleDropZoneDragOver(event: DragEvent) {
        event.preventDefault();
        isDraggingOver = true;
    }

    function handleDropZoneDragLeave() {
        isDraggingOver = false;
    }

    function handleDropZoneDrop(event: DragEvent) {
        event.preventDefault();
        isDraggingOver = false;

        if (event.dataTransfer?.files) {
            addFiles(Array.from(event.dataTransfer.files));
        }
    }

    async function handleSubmit() {
        submitError = "";
        submitSuccess = false;

        // Validation
        if (!name.trim()) {
            submitError = "Name is required";
            return;
        }
        if (!cleanInstagramUsername) {
            submitError = "Instagram username is required";
            return;
        }
        if (instagramValid === false) {
            submitError = "Please enter a valid Instagram username";
            return;
        }
        if (!caption.trim()) {
            submitError = "Caption is required";
            return;
        }
        if (isOverLimit) {
            submitError = "Caption exceeds 2,200 characters";
            return;
        }
        if (uploadFiles.length === 0) {
            submitError = "At least one image or video is required";
            return;
        }

        isSubmitting = true;

        try {
            // Create FormData
            const formData = new FormData();
            formData.append("andrewId", andrewId);
            formData.append("name", name.trim());
            formData.append("instagramUsername", cleanInstagramUsername);
            formData.append("caption", caption.trim());

            // Add files in order
            uploadFiles
                .sort((a, b) => a.order - b.order)
                .forEach((uf, index) => {
                    formData.append(`file_${index}`, uf.file);
                });

            const accessToken = getAccessToken();

            const response = await fetch("/api/posts/submit", {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${accessToken}`,
                },
                body: formData,
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || "Failed to submit post");
            }

            // Success!
            submitSuccess = true;
            caption = "";
            instagramUsername = "";
            instagramValid = null;
            uploadFiles.forEach((f) => URL.revokeObjectURL(f.preview));
            uploadFiles = [];
        } catch (e) {
            submitError =
                e instanceof Error ? e.message : "Failed to submit post";
        } finally {
            isSubmitting = false;
        }
    }

    function handleLogout() {
        logout();
    }
</script>

<div class="min-h-screen bg-gray-100">
    <!-- Header -->
    <header class="bg-white shadow-sm sticky top-0 z-10">
        <div
            class="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between"
        >
            <h1 class="text-2xl font-bold text-red-600">CMU Insta</h1>
            <div class="flex items-center gap-4">
                <span class="text-gray-600 text-sm">
                    Signed in as <span class="font-medium">{andrewId}</span>
                </span>
                <button
                    on:click={handleLogout}
                    class="px-4 py-2 text-sm text-gray-600 hover:text-red-600 transition-colors"
                >
                    Sign Out
                </button>
            </div>
        </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-2xl mx-auto px-4 py-8">
        <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
            <h2 class="text-2xl font-semibold text-gray-800 mb-6">
                Create a Post
            </h2>

            {#if submitSuccess}
                <div
                    class="mb-6 p-4 bg-green-50 border border-green-200 rounded-xl"
                >
                    <div class="flex items-center gap-3">
                        <svg
                            class="w-6 h-6 text-green-500"
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
                        <div>
                            <p class="text-green-800 font-medium">
                                Post submitted successfully!
                            </p>
                            <p class="text-green-600 text-sm">
                                Your post is pending admin approval.
                            </p>
                        </div>
                    </div>
                </div>
            {/if}

            {#if submitError}
                <div
                    class="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl"
                >
                    <p class="text-red-700">{submitError}</p>
                </div>
            {/if}

            <form on:submit|preventDefault={handleSubmit} class="space-y-6">
                <!-- Name Field -->
                <div>
                    <label
                        for="name"
                        class="block text-sm font-medium text-gray-700 mb-2"
                    >
                        Name
                    </label>
                    <input
                        type="text"
                        id="name"
                        bind:value={name}
                        placeholder="Your name"
                        class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-shadow"
                    />
                </div>

                <!-- Andrew ID Field (Read-only) -->
                <div>
                    <label
                        for="andrewId"
                        class="block text-sm font-medium text-gray-700 mb-2"
                    >
                        Andrew ID
                    </label>
                    <input
                        type="text"
                        id="andrewId"
                        value={andrewId}
                        readonly
                        class="w-full px-4 py-3 border border-gray-200 rounded-xl bg-gray-50 text-gray-600 cursor-not-allowed"
                    />
                </div>

                <!-- Instagram Username Field -->
                <div>
                    <label
                        for="instagramUsername"
                        class="block text-sm font-medium text-gray-700 mb-2"
                    >
                        Instagram Username
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400"
                            >@</span
                        >
                        <input
                            type="text"
                            id="instagramUsername"
                            bind:value={instagramUsername}
                            on:input={handleInstagramInput}
                            placeholder="your_instagram"
                            class="w-full pl-8 pr-12 py-3 border rounded-xl focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-shadow
                                   {instagramValid === false
                                ? 'border-red-500 bg-red-50'
                                : 'border-gray-300'}"
                        />
                        <div class="absolute right-4 top-1/2 -translate-y-1/2">
                            {#if instagramValid === false}
                                <svg
                                    class="h-5 w-5 text-red-500"
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
                            {/if}
                        </div>
                    </div>
                    {#if instagramError}
                        <p class="mt-1 text-sm text-red-600">
                            {instagramError}
                        </p>
                    {:else}
                        <p class="mt-1 text-sm text-gray-500">
                            Enter your Instagram username so we can tag you
                        </p>
                    {/if}

                    <!-- Instagram Profile Preview Card -->
                    {#if cleanInstagramUsername && instagramValid === true}
                        <div
                            class="mt-3 p-4 border border-gray-200 rounded-xl bg-gradient-to-r from-purple-50 via-pink-50 to-orange-50"
                        >
                            <div class="flex items-center gap-4">
                                <!-- Instagram Icon -->
                                <div
                                    class="w-14 h-14 rounded-full bg-gradient-to-tr from-yellow-400 via-pink-500 to-purple-600 p-[2px] flex-shrink-0"
                                >
                                    <div
                                        class="w-full h-full rounded-full bg-white flex items-center justify-center"
                                    >
                                        <svg
                                            class="w-7 h-7"
                                            viewBox="0 0 24 24"
                                            fill="url(#instagram-gradient)"
                                        >
                                            <defs>
                                                <linearGradient
                                                    id="instagram-gradient"
                                                    x1="0%"
                                                    y1="100%"
                                                    x2="100%"
                                                    y2="0%"
                                                >
                                                    <stop
                                                        offset="0%"
                                                        style="stop-color:#FFDC80"
                                                    />
                                                    <stop
                                                        offset="25%"
                                                        style="stop-color:#F77737"
                                                    />
                                                    <stop
                                                        offset="50%"
                                                        style="stop-color:#E1306C"
                                                    />
                                                    <stop
                                                        offset="75%"
                                                        style="stop-color:#C13584"
                                                    />
                                                    <stop
                                                        offset="100%"
                                                        style="stop-color:#833AB4"
                                                    />
                                                </linearGradient>
                                            </defs>
                                            <path
                                                d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"
                                            />
                                        </svg>
                                    </div>
                                </div>

                                <!-- Profile Info -->
                                <div class="flex-1 min-w-0">
                                    <p
                                        class="font-semibold text-gray-900 truncate"
                                    >
                                        @{cleanInstagramUsername}
                                    </p>
                                    <p class="text-sm text-gray-500">
                                        Please verify this is your account
                                    </p>
                                </div>

                                <!-- View Profile Button -->
                                <button
                                    type="button"
                                    on:click={openInstagramProfile}
                                    class="flex-shrink-0 px-4 py-2 bg-gradient-to-r from-purple-500 via-pink-500 to-orange-500 text-white text-sm font-medium rounded-lg hover:opacity-90 transition-opacity flex items-center gap-2"
                                >
                                    <svg
                                        class="w-4 h-4"
                                        fill="none"
                                        stroke="currentColor"
                                        viewBox="0 0 24 24"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                                        ></path>
                                    </svg>
                                    Verify
                                </button>
                            </div>
                        </div>
                    {/if}
                </div>

                <!-- Caption Field -->
                <div>
                    <label
                        for="caption"
                        class="block text-sm font-medium text-gray-700 mb-2"
                    >
                        Caption
                    </label>
                    <textarea
                        id="caption"
                        bind:value={caption}
                        placeholder="Write a caption for your post..."
                        rows="4"
                        class="w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-shadow resize-none
                               {isOverLimit
                            ? 'border-red-500 bg-red-50'
                            : 'border-gray-300'}"
                    ></textarea>
                    <div class="flex justify-end mt-1">
                        <span
                            class="text-sm {isOverLimit
                                ? 'text-red-600 font-medium'
                                : 'text-gray-500'}"
                        >
                            {charCount}/2200
                        </span>
                    </div>
                </div>

                <!-- File Upload -->
                <div>
                    <label
                        for="file-upload"
                        class="block text-sm font-medium text-gray-700 mb-2"
                    >
                        Photos & Videos
                        <span class="text-gray-400 font-normal"
                            >({uploadFiles.length}/10)</span
                        >
                    </label>

                    <!-- Drop Zone -->
                    <div
                        bind:this={dropTarget}
                        on:dragover={handleDropZoneDragOver}
                        on:dragleave={handleDropZoneDragLeave}
                        on:drop={handleDropZoneDrop}
                        class="border-2 border-dashed rounded-xl p-6 text-center transition-colors cursor-pointer
                               {isDraggingOver
                            ? 'border-red-500 bg-red-50'
                            : 'border-gray-300 hover:border-gray-400'}"
                        on:click={() => fileInput.click()}
                        on:keypress={(e) =>
                            e.key === "Enter" && fileInput.click()}
                        role="button"
                        tabindex="0"
                    >
                        <input
                            id="file-upload"
                            bind:this={fileInput}
                            type="file"
                            accept="image/*,video/*"
                            multiple
                            class="hidden"
                            on:change={handleFileSelect}
                            disabled={uploadFiles.length >= 10}
                        />
                        <svg
                            class="w-10 h-10 mx-auto text-gray-400 mb-3"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                            ></path>
                        </svg>
                        <p class="text-gray-600 mb-1">
                            {#if uploadFiles.length >= 10}
                                Maximum files reached
                            {:else}
                                Drag and drop or click to upload
                            {/if}
                        </p>
                        <p class="text-gray-400 text-sm">
                            Images and videos • Max 10 files
                        </p>
                    </div>

                    <!-- File Previews -->
                    {#if uploadFiles.length > 0}
                        <div class="mt-4">
                            <p class="text-xs text-gray-500 mb-2">
                                Drag to reorder. First image/video will be the
                                cover.
                            </p>
                            <div class="grid grid-cols-5 gap-2">
                                {#each uploadFiles as uf (uf.id)}
                                    <div
                                        draggable="true"
                                        on:dragstart={(e) =>
                                            handleDragStart(e, uf)}
                                        on:dragover={handleDragOver}
                                        on:dragenter={(e) =>
                                            handleDragEnter(e, uf)}
                                        on:dragend={handleDragEnd}
                                        class="relative aspect-square rounded-lg overflow-hidden bg-gray-100 cursor-move group
                                               {draggedItem?.id === uf.id
                                            ? 'opacity-50'
                                            : ''}"
                                        role="listitem"
                                    >
                                        {#if uf.file.type.startsWith("video/")}
                                            <video
                                                src={uf.preview}
                                                class="w-full h-full object-cover"
                                                muted
                                            >
                                                <track kind="captions" />
                                            </video>
                                            <div
                                                class="absolute inset-0 flex items-center justify-center bg-black/30"
                                            >
                                                <svg
                                                    class="w-8 h-8 text-white"
                                                    fill="currentColor"
                                                    viewBox="0 0 24 24"
                                                >
                                                    <path d="M8 5v14l11-7z"
                                                    ></path>
                                                </svg>
                                            </div>
                                        {:else}
                                            <img
                                                src={uf.preview}
                                                alt="Upload preview"
                                                class="w-full h-full object-cover"
                                            />
                                        {/if}

                                        <!-- Order badge -->
                                        <div
                                            class="absolute top-1 left-1 w-5 h-5 bg-black/70 rounded-full flex items-center justify-center"
                                        >
                                            <span
                                                class="text-white text-xs font-medium"
                                                >{uf.order}</span
                                            >
                                        </div>

                                        <!-- Remove button -->
                                        <button
                                            type="button"
                                            on:click|stopPropagation={() =>
                                                removeFile(uf.id)}
                                            class="absolute top-1 right-1 w-6 h-6 bg-red-500 hover:bg-red-600 rounded-full flex items-center justify-center
                                                   opacity-0 group-hover:opacity-100 transition-opacity"
                                        >
                                            <svg
                                                class="w-4 h-4 text-white"
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
                                        </button>
                                    </div>
                                {/each}
                            </div>
                        </div>
                    {/if}
                </div>

                <!-- Submit Button -->
                <button
                    type="submit"
                    disabled={isSubmitting || isOverLimit}
                    class="w-full py-4 px-6 bg-red-600 hover:bg-red-700 disabled:bg-red-400
                           text-white font-semibold rounded-xl transition-all duration-200
                           flex items-center justify-center gap-2 shadow-lg hover:shadow-xl
                           disabled:cursor-not-allowed"
                >
                    {#if isSubmitting}
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
                        <span>Submitting...</span>
                    {:else}
                        <svg
                            class="w-5 h-5"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                            ></path>
                        </svg>
                        <span>Submit Post</span>
                    {/if}
                </button>
            </form>
        </div>
    </main>
</div>
