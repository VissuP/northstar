import {Component} from "@angular/core";
import {HashLocationStrategy, Location, LocationStrategy, NgClass} from "@angular/common";
import {Router, NavigationStart, Event} from "@angular/router";

import {BreadcrumbService} from "./breadcrumb.service";

@Component({
    providers: [
        { provide: LocationStrategy, useClass: HashLocationStrategy },
    ],
    selector: "breadcrumb",
    styleUrls: ["breadcrumb.component.css"],
    templateUrl: "breadcrumb.component.html",
})
export class BreadcrumbComponent {

  private router: Router;
  private location: Location;
  private breadcrumbService: BreadcrumbService;

  // The route url used to lookup display name. Note
  // that by default we use the root level route,
  private url: string = "/";

  constructor(router: Router, location: Location, breadcrumbService: BreadcrumbService) {
    this.router = router;
    this.location = location;
    this.breadcrumbService = breadcrumbService;

    this.router.events.filter((e: Event) => e instanceof NavigationStart)
    .subscribe((event: NavigationStart) => {
      if ((event.constructor as any ).name === "NavigationEnd") {
        this.url = event.url;
      }
    });
  }

  // Returns the display name of the current route.
  public getBreadcrumb(): String {
    return this.breadcrumbService.getRouteDisplayName(this.url);
  }
}
