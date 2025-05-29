import { ChartAreaInteractive } from "@/components/chart-area-interactive";
import { DataTable } from "@/components/data-table";
import { MonitorsCards } from "@/components/monitors-cards";
import { SiteHeader } from "@/components/site-header";

import data from "../../data.json";

export default async function Page({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return (
    <>
      <SiteHeader title={`Monitor: ${id}`} />
      <div className="flex flex-1 flex-col">
        <div className="@container/main flex flex-1 flex-col gap-4">
          <div className="flex flex-col gap-6 py-6">
            {/* Monitor Statistics */}
            <MonitorsCards />

            {/* Monitor Results Chart */}
            <div className="px-4 lg:px-6">
              <ChartAreaInteractive />
            </div>

            {/* Monitor Results Table */}
            <DataTable data={data} />
          </div>
        </div>
      </div>
    </>
  );
}
