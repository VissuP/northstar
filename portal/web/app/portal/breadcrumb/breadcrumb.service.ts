import { Injectable } from "@angular/core";

@Injectable()
export class BreadcrumbService {

    // Defines the map of display names for routes.
    private routeDisplayNames: Map<string, string> = new Map<string, string>();

    // Specifies a displayName for a route. Note that route should be
    // the full url as in the same url used to call router.navigate().
    public addRouteDisplayName(route: string, displayName: string): void {
        this.routeDisplayNames[route] = displayName;
    }

    // Returns a displayName for a route.
    public getRouteDisplayName(route: string): string {
        let displayName = this.routeDisplayNames[route];
        if (!displayName) {
            displayName = route.substr(1, route.length);
        }

        return displayName;
    }
}