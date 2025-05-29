import { ChartAreaInteractive } from "@/components/chart-area-interactive";
import { DataTable } from "@/components/data-table";
import { SectionCards } from "@/components/section-cards";
import { SiteHeader } from "@/components/site-header";

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
        <div className="@container/main flex flex-1 flex-col gap-2">
          <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
            {/* Basic Stats */}
            <SectionCards />

            {/* (Live) Monitor Results Chart */}
            <div className="px-4 lg:px-6">
              <ChartAreaInteractive />
            </div>

          </div>
        </div>
      </div>
    </>
  );
}
