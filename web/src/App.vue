<script setup lang="ts">
import { onMounted, ref, type Ref } from 'vue';
import type { FeatureToggle } from './core/types';
import DialogCreate from './components/DialogCreate.vue';
import DialogEdit from './components/DialogEdit.vue';
import Toggle from './components/Toggle.vue';

const showEditDialog = ref(false)
const showCreateDialog = ref(false)
const featureToggles: Ref<Array<FeatureToggle>> = ref([])
const dialogModel = ref<FeatureToggle>({
  featureId: '',
  description: '',
  enabled: false,
})

const handleClick = (ev: Event, featureToggle: FeatureToggle) => {
  const cb = ev.target as HTMLInputElement

  const updatedToggles = featureToggles.value.map(ft => ({
    ...ft,
    enabled: ft.featureId === featureToggle.featureId ? cb.checked : ft.enabled,
  }))

  featureToggles.value = updatedToggles;
}

const getToggles = async (): Promise<void> => {
  const response = await fetch('/api/v1/features')
  if (response.status !== 200) {
    console.error(`could not fetch feature toggles, status code ${response.status}`)
    return
  }

  featureToggles.value = await response.json()
}

const updateToggle = async (toggle: FeatureToggle): Promise<Response> => {
  return fetch('/api/v1/features', {
    method: 'POST',
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(toggle),
  })
}

const apply = async (): Promise<void> => {
  await Promise.all(featureToggles.value.map(toggle => updateToggle(toggle)))
  return getToggles()
}

const createOrUpdate = async () => {
  await updateToggle(dialogModel.value)
  await getToggles()
}

const remove = async (toggle: FeatureToggle) => {
  await fetch(`/api/v1/features/${toggle.featureId}`, {
    method: 'DELETE',
  })
  return getToggles()
}

const openCreate = () => {
  dialogModel.value = {
    featureId: '',
    description: '',
    enabled: false,
  }
  showCreateDialog.value = true
}

const openEdit = (selected: FeatureToggle) => {
  dialogModel.value = { ...selected }
  showEditDialog.value = true
}

onMounted(getToggles)
</script>

<template>
  <main class="flex flex-col items-center min-h-screen dark:bg-slate-900 p-16">
    <div class="w-[800px]">
      <h1 class="text-xl font-bold text-gray-200 mb-8">Feature Toggles</h1>
      <ul class="w-full">
        <li class="flex justify-between items-center mb-5 gap-4" v-for="featureToggle in featureToggles">
          <div>
            <Toggle :toggle="featureToggle" @change="enabled => featureToggle.enabled = enabled" />
            <p class="dark:text-gray-500">{{ featureToggle.description }}</p>
          </div>
          <div class="flex gap-2">
            <button class="bg-blue-600 rounded px-3 py-1 text-sm text-gray-200" @click="openEdit(featureToggle)">Edit</button>
            <button class="bg-red-400 rounded px-3 py-1 text-sm text-gray-200" @click="remove(featureToggle)">X</button>
          </div>
        </li>
      </ul>
      <div class="flex gap-2">
        <button class="bg-blue-600 rounded px-3 py-1 text-sm text-gray-200" @click="apply">Apply</button>
        <button class="bg-blue-600 rounded px-3 py-1 text-sm text-gray-200" @click="openCreate">Create new</button>
      </div>

      <DialogCreate v-model:show="showCreateDialog" v-model:toggle="dialogModel" @apply="createOrUpdate" />
      <DialogEdit v-model:show="showEditDialog" v-model:toggle="dialogModel" @apply="createOrUpdate" />
    </div>
  </main>
</template>