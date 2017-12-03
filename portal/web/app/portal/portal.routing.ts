import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { HelpComponent } from "./help/help.component";
import { BreadcrumbService } from "./breadcrumb/breadcrumb.service";
import { EditorComponent } from "./notebook/editor/editor.component";
import { NotebookComponent } from "./notebook/notebook.component";
import { PortfolioComponent } from "./portfolio/portfolio.component";
import { PortfolioResolver } from "./portfolio/file/file.resolver";
import { FileComponent } from "./portfolio/file/file.component";
import { TransformationComponent } from "./transformation/transformation.component";
import { TransformationResults } from "./transformation/results/results.component";
import { TransformationResultsResolver } from "./transformation/results/results.resolver";
import { NotebookResolver } from "./notebook/editor/editor.resolver";
import { ErrorComponent } from "./error/error.component";
import { AuthGuard } from './../shared/auth/auth-guard.service';
import { PortalComponent } from './../portal/portal.component';
// Defines the portal (child) routes.
const portalRoutes: Routes = [
    {
        path: "portal",
        component: PortalComponent,
        canActivate: [AuthGuard],
        canActivateChild: [AuthGuard],

        children: [
            {
                path: '',
                component: DashboardComponent,
            },
            {
                path: 'dashboard',
                component: DashboardComponent,
            },
            {
                path: 'notebooks',
                component: NotebookComponent,
            },
            {
                path: 'editor',
                component: EditorComponent,
            },
            {
                path: 'editor/:id',
                component: EditorComponent,
                resolve: {
                    notebook: NotebookResolver,
                },
            },
            {
                path: 'portfolios',
                component: PortfolioComponent,
            },
            {
                path: 'portfolios/:id',
                component: FileComponent,
                resolve: {
                    portfolio: PortfolioResolver,
                },
            },
            {
                path: 'transformations',
                component: TransformationComponent,
            },
            {
                path: 'transformations/:id/results',
                component: TransformationResults,
                resolve: {
                    history: TransformationResultsResolver,
                },
            },
            {
                path: 'help',
                component: HelpComponent,
            },
            {
                path: 'error',
                component: ErrorComponent,
            }
        ]
    },
];

@NgModule({
    exports: [
        RouterModule,
    ],
    imports: [
        RouterModule.forChild(portalRoutes),
    ],
    providers: [
        BreadcrumbService,
    ],
})
export class PortalRoutingModule {
    constructor(breadcrumbService: BreadcrumbService) {
        // Setup the root level routes breadcrumb display name.
        breadcrumbService.addRouteDisplayName('', 'Home / Dashboard');
        breadcrumbService.addRouteDisplayName('/', 'Home / Dashboard');
        breadcrumbService.addRouteDisplayName('/portal', 'Home / Dashboard');

        // For every entry in the portal routes, setup the breadcrumb
        // display name. Note that we assume childs will be under
        // "/portal".
        for (let childRoute of portalRoutes) {
            if (childRoute.path !== "") {
                let route = "/portal/" + childRoute.path;
                let displayName = "Home / " + childRoute.path;
                breadcrumbService.addRouteDisplayName(route, displayName);
            }
        }
    }
}
