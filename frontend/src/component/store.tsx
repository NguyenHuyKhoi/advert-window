import { create } from "zustand";

export interface IDeviceInfo {
  device_id: string;
  device_model: string;
  device_name: string;
  os_version: string;
}

export interface IDeviceAdvertStatus {
  id?: string;
  playing?: boolean;
  volume?: number;
}

export interface IAdvert {
  id: string;
  name: string;
  audio_url: string;
}

interface AppState {
  deviceInfo: IDeviceInfo | null;
  targetAdvert: IAdvert | null;
  advertStatus: IDeviceAdvertStatus | null;
  deviceFields: any[];
  isInfoModalOpen: boolean; // State quản lý modal
  setDeviceInfo: (info: IDeviceInfo) => void;
  setDeviceFields: (fields: any[]) => void;
  setTargetAdvert: (advert: IAdvert) => void;
  updateAdvertStatus: (status: Partial<IDeviceAdvertStatus>) => void;
  setModalOpen: (isOpen: boolean) => void; // Hàm trigger modal
}

export const useAppStore = create<AppState>((set) => ({
  deviceInfo: null,
  targetAdvert: null,
  advertStatus: null,
  deviceFields: [],
  isInfoModalOpen: false,
  setDeviceInfo: (info) => set({ deviceInfo: info }),
  setDeviceFields: (fields) => set({ deviceFields: fields }),
  setTargetAdvert: (advert) => {
    set({
      targetAdvert: advert,
    });
  },
  updateAdvertStatus: (status) =>
    set((state) => ({
      advertStatus: { ...(state.advertStatus ?? {}), ...status },
    })),
  setModalOpen: (isOpen) => set({ isInfoModalOpen: isOpen }),
}));
