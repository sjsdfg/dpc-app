package gov.cms.dpc.fhir.converters;

import org.hl7.fhir.r4.model.Base;

public interface FHIRConverter<R extends Base, C> {

    C fromFHIR(FHIREntityConverter converter, R resource);

    R toFHIR(FHIREntityConverter converter, C entity);

    Class<R> getFHIRResource();

    Class<C> getJavaClass();
}
