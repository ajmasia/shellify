import { makeAutoObservable } from 'mobx'

export class UIStore {
  currentProjectId: string | null = null
  currentProjectName: string | null = null

  constructor() {
    makeAutoObservable(this)
  }

  setCurrentProject(id: string | null, name: string | null): void {
    this.currentProjectId = id
    this.currentProjectName = name
  }

  clearProject(): void {
    this.currentProjectId = null
    this.currentProjectName = null
  }
}
