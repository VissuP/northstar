import {BootstrapModalModule} from "angular2-modal/plugins/bootstrap";
import {NgModule, ModuleWithProviders} from "@angular/core";
import {CommonModule} from "@angular/common";
import {FormsModule} from "@angular/forms";
import {HttpModule} from "@angular/http";
import {ModalModule} from "angular2-modal";
import {CodeModule} from "ngx-vz-cell";

import {TooltipModule} from "ng2-tooltip";
import {Ng2AutoCompleteModule} from "ng2-auto-complete";
import "../rxjs-extensions";
import {NotebookService} from "../shared/services/notebook.service";
import {UserService} from "../shared/services/user.service";
import {EventService} from "ngx-vz-websocket";
import {TransformationService} from "./transformation/transformation.service";
import {TemplateService} from "../shared/services/template.service";

import {AlertService} from "../shared/services/alerts.service";
import { NgxDatatableModule }               from '@swimlane/ngx-datatable';

/* Portal Components */
import {PortalRoutingModule} from "./portal.routing";
import {PortalComponent} from "./portal.component";
import {DashboardComponent} from "./dashboard/dashboard.component";
import {BreadcrumbComponent} from "./breadcrumb/breadcrumb.component";
import {HelpComponent} from "./help/help.component";

/* Notebook Components */
import {NotebookComponent} from "./notebook/notebook.component";
import {EditorComponent} from "./notebook/editor/editor.component";
import {Ng2CompleterModule} from "ng2-completer";
import {StaticCellComponent} from "./notebook/editor/cell/static/static.component";
import {PermissionsModal} from "./notebook/editor/permissions/permissions.component";

/* Transformation Components */
import {TransformationComponent} from "./transformation/transformation.component";
import {TransformationModal} from "./transformation/editor/editor.component";
import {LanguageService} from "../shared/services/language.service";
import {TransformationResults} from "./transformation/results/results.component";
import {CodemirrorModule} from "ng2-codemirror";

/* Portfolio Components */
import {PortfolioComponent} from "./portfolio/portfolio.component";
import {PortfolioService} from "../shared/services/portfolio.service";
import {FileComponent} from "./portfolio/file/file.component";

/* Shared Components */
import {AlertsComponent} from "../shared/components/alerts/alerts.component";
import {NotebookResolver} from "./notebook/editor/editor.resolver";
import {PortfolioResolver} from "./portfolio/file/file.resolver";
import {GlobalErrorHandlerProvider} from "../shared/resources/error";
import {ErrorComponent} from "./error/error.component";
import {TransformationResultsResolver} from "./transformation/results/results.resolver";

@NgModule({
    imports: [
        CommonModule,
        FormsModule,
        HttpModule,
        PortalRoutingModule,
        CodeModule,
        CodemirrorModule,
        ModalModule.forRoot(),
        BootstrapModalModule,
        TooltipModule,
        Ng2CompleterModule,
        Ng2AutoCompleteModule,
        NgxDatatableModule
    ],
    declarations: [
        /* Portal */
        PortalComponent,
        DashboardComponent,
        HelpComponent,
        EditorComponent,
        NotebookComponent,
        BreadcrumbComponent,
        StaticCellComponent,
        PermissionsModal,
        TransformationResults,
        PortfolioComponent,
        FileComponent,

        TransformationComponent,
        TransformationModal,
        AlertsComponent,
        ErrorComponent,
    ],
    providers: [
        UserService,
        NotebookService,
        EventService,
        PortfolioService,
        TransformationService,
        LanguageService,
        AlertService,
        TemplateService,
        NotebookResolver,
        PortfolioResolver,
        TransformationResultsResolver,
        GlobalErrorHandlerProvider,
    ],
    entryComponents: [
        PermissionsModal,
        TransformationModal,
    ],
})
export class PortalModule {
      static forRoot(): ModuleWithProviders {
    return {
      ngModule: PortalModule
    };
  }
}
