import { IconTrendingDown, IconTrendingUp, IconClock, IconAlertCircle, IconCheck, IconActivity } from "@tabler/icons-react";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardAction,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Link from "next/link";
import { cn } from "@/lib/utils"; // Utility for conditional classnames

interface Monitor {
  Id: string;
  Url: string;
  Active: boolean;
  CreatedAt: string;
  Uptime?: number; // Add optional mock metrics
  AvgResponseTimeMs?: number;
  MonthlyRequests?: number;
  TrendingUp?: boolean;
}

export function MonitorCards({ data }: { data: Monitor[] }) {
  return (
    <div className="grid grid-cols-1 gap-4 px-4 sm:grid-cols-2 xl:grid-cols-3">
      {data.map((monitor) => (
        <Link href={`/dashboard/monitors/${monitor.Id}`} key={monitor.Id}>
          <Card className="hover:shadow-xl transition-shadow duration-200 border border-border rounded-2xl p-4 bg-card">
            <CardHeader className="pb-2">
              <div className="flex justify-between items-start">
                <div>
                  <CardDescription className="text-xs text-muted-foreground mb-1">
                    ID: {monitor.Id.slice(0, 8)}...
                  </CardDescription>
                  <CardTitle className="text-lg font-bold break-words">
                    {monitor.Url}
                  </CardTitle>
                </div>
                <CardAction>
                  <Badge variant={monitor.Active ? "default" : "secondary"}>
                    {monitor.Active ? "Active" : "Inactive"}
                  </Badge>
                </CardAction>
              </div>
            </CardHeader>

            <CardFooter className="pt-2 flex flex-col gap-3 text-sm">
              <div
                className={cn(
                  "flex items-center gap-1.5 font-medium",
                  monitor.TrendingUp ? "text-green-600" : "text-red-600"
                )}
              >
                {monitor.TrendingUp ? (
                  <>
                    Trending up this month <IconTrendingUp className="size-4" />
                  </>
                ) : (
                  <>
                    Trending down <IconTrendingDown className="size-4" />
                  </>
                )}
              </div>

              <div className="grid grid-cols-3 gap-3 text-xs text-muted-foreground">
                <div>
                  <div className="font-medium text-foreground">
                    {monitor.Uptime ?? "99.9"}%
                  </div>
                  <div>Uptime</div>
                </div>
                <div>
                  <div className="font-medium text-foreground">
                    {monitor.AvgResponseTimeMs ?? 180}ms
                  </div>
                  <div>Avg Response</div>
                </div>
                <div>
                  <div className="font-medium text-foreground">
                    {monitor.MonthlyRequests ?? 2500}
                  </div>
                  <div>Monthly Hits</div>
                </div>
              </div>

              <div className="text-xs text-muted-foreground mt-1">
                Created: {new Date(monitor.CreatedAt).toLocaleString()}
              </div>
            </CardFooter>
          </Card>
        </Link>
      ))}
    </div>
  );
}

export function MonitorsCards() {
  return (
    <div className="*:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card grid grid-cols-1 gap-4 px-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:shadow-xs lg:px-6 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Uptime</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            99.98%
          </CardTitle>
          <CardAction>
            <Badge variant="outline" className="bg-green-500/10 text-green-500">
              <IconCheck className="size-4" />
              Operational
            </Badge>
          </CardAction>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Last 30 days <IconClock className="size-4" />
          </div>
          <div className="text-muted-foreground">
            Only 0.02% downtime
          </div>
        </CardFooter>
      </Card>
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Response Time</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            245ms
          </CardTitle>
          <CardAction>
            <Badge variant="outline" className="bg-yellow-500/10 text-yellow-500">
              <IconActivity className="size-4" />
              Good
            </Badge>
          </CardAction>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Average response time <IconActivity className="size-4" />
          </div>
          <div className="text-muted-foreground">
            Target: {"<"} 300ms
          </div>
        </CardFooter>
      </Card>
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Alerts</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            3
          </CardTitle>
          <CardAction>
            <Badge variant="outline" className="bg-orange-500/10 text-orange-500">
              <IconAlertCircle className="size-4" />
              Active
            </Badge>
          </CardAction>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Last 24 hours <IconAlertCircle className="size-4" />
          </div>
          <div className="text-muted-foreground">
            2 critical, 1 warning
          </div>
        </CardFooter>
      </Card>
      <Card className="@container/card">
        <CardHeader>
          <CardDescription>Checks</CardDescription>
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            1,234
          </CardTitle>
          <CardAction>
            <Badge variant="outline" className="bg-blue-500/10 text-blue-500">
              <IconActivity className="size-4" />
              Today
            </Badge>
          </CardAction>
        </CardHeader>
        <CardFooter className="flex-col items-start gap-1.5 text-sm">
          <div className="line-clamp-1 flex gap-2 font-medium">
            Total checks performed <IconActivity className="size-4" />
          </div>
          <div className="text-muted-foreground">
            Every 5 minutes
          </div>
        </CardFooter>
      </Card>
    </div>
  )
}
