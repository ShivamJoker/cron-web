import Alpine from "alpinejs";
import Wretch from "wretch";
import cronparse from "cronstrue";
import { Notyf } from "notyf";
import "notyf/notyf.min.css";
const notyf = new Notyf({ position: { y: "top", x: "center" } });

const API = Wretch("./crontab", {}).resolve((res) => res.json<CronItem[]>());

type CronItem = {
  Cron: string;
  Job: string;
  editing: boolean;
  explanation: string;
  isValid: boolean;
};

const fetchLatestCrontabs = async () => {
  const res = await API.get();
  const updatedRes = res.map((c) => ({
    ...c,
    explanation: cronparse.toString(c.Cron),
    isValid: true,
    editing: false,
  }));
  return updatedRes;
};

Alpine.data("cronjobs", () => ({
  items: [] as CronItem[],

  refreshing: false,

  async refreshCrontabs() {
    this.refreshing = true;
    this.items = await fetchLatestCrontabs();
    this.refreshing = false;
  },

  toggleEdit(idx: number) {
    this.items[idx].editing = true;
  },

  async saveEdit(idx: number) {
    if (!this.items[idx].isValid) return;
    this.items[idx].editing = false;
    await API.put(this.items);
    notyf.success("Updated CRON tab");
  },

  addNewCrontab() {
    const sampleCron = "* * * * *";
    this.items.push({
      editing: true,
      isValid: true,
      Cron: sampleCron,
      Job: "echo echo",
      explanation: cronparse.toString(sampleCron),
    });
  },
  async deleteCrontab(deleteIdx: number) {
    this.items = this.items.filter((_, idx) => idx !== deleteIdx);
    await API.put(this.items);
    notyf.error("Deleted CRON tab");
  },
  async init() {
    this.items = await fetchLatestCrontabs();

    this.$watch("items", () => {
      this.items.forEach((item, idx) => {
        if (!item.editing) return;
        try {
          this.items[idx].explanation = cronparse.toString(item.Cron);
          this.items[idx].isValid = true;
        } catch (error) {
          this.items[idx].isValid = false;
          this.items[idx].explanation = "Invalid CRON expression";
        }
      });
    });
  },
}));

Alpine.start();
