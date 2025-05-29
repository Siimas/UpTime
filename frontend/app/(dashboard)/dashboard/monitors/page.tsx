import { MonitorCards } from "@/components/monitors-cards";
import { SiteHeader } from "@/components/site-header";

const monitors = [
  {
    Id: "db9e1955-0f53-4832-96fb-b2479331b2a3",
    Url: "https://httpstat.us/200",
    Active: true,
    IntervalSeconds: 10,
    CreatedAt: "2025-05-27 18:41:44.97926+00"
  },
  {
    Id: "d7c77caf-8a01-4fef-b88c-c341484c800b",
    Url: "https://httpstat.us/400",
    Active: false,
    IntervalSeconds: 35,
    CreatedAt: "2025-05-28 19:26:03.505866+00"
  },
]

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
            
            <MonitorCards data={monitors} />

          </div>
        </div>
      </div>
    </>
  );
}
