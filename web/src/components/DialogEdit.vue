<script setup lang="ts">
import type { FeatureToggle } from '@/core/types';
import { ref, watch } from 'vue';

const emit = defineEmits<{
    (e: 'update:show', show: boolean): void
    (e: 'update:toggle', toggle: FeatureToggle): void
    (e: 'cancel'): void
    (e: 'apply'): void
}>()

const props = defineProps<{
    show?: boolean
    toggle: FeatureToggle
}>()

const show = ref(props.show)
const toggle = ref(props.toggle)

watch(() => props.show, (newValue) => {
    show.value = newValue
})
watch(() => props.toggle, (newValue) => {
    toggle.value = newValue
})

const apply = () => {
    show.value = false

    emit('apply')
    emit('update:toggle', toggle.value)
    emit('update:show', show.value)
}

const cancel = () => {
    show.value = false
    emit('update:show', show.value)
}
</script>

<template>
    <div v-show="show" class="absolute flex top-0 left-0 w-screen h-screen bg-black/60 justify-center items-center">
        <div class="absolute top-0 left-0 w-full h-full z-0" @click="cancel"></div>
        <div class="relative bg-slate-900 rounded-md p-6 shadow-xl w-[500px] z-1">
            <h2 class="text-lg font-medium text-gray-200 mb-4">Create feature toggle</h2>
            <div class="mb-4">
                <label for="id" class="text-gray-500 block text-sm">Identifier</label>
                <input disabled v-model="toggle.featureId" name="id" type="text" class="bg-slate-700 rounded text-md text-gray-200 p-1 mb-2 w-full">
                <label for="description" class="text-gray-500 block text-sm">Description</label>
                <textarea v-model="toggle.description" name="description" class="bg-slate-700 rounded text-md text-gray-200 p-1 mb-2 w-full" />
            </div>
            <div class="flex flex-row justify-end gap-2">
                <button class="bg-blue-600 rounded px-3 py-1 text-sm text-gray-200" @click="apply">Save</button>
                <button class="bg-blue-600 rounded px-3 py-1 text-sm text-gray-200" @click="cancel">Cancel</button>
            </div>
        </div>
    </div>
</template>