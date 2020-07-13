CREATE OR REPLACE FUNCTION getImporteEnLetras(importe numeric)
    RETURNS character varying AS
$BODY$

BEGIN

    importe = importe::numeric(16,2);


    IF split_part(importe::VARCHAR,'.',2) = '00' THEN
        IF split_part(importe::VARCHAR, '.',1) = '1' THEN
            RETURN concat_ws(' ','Pesos uno')::VARCHAR;
        ELSE
            RETURN concat_ws(' ','Pesos', lower(replace(fu_numero_letras(coalesce(split_part(importe::VARCHAR, '.',1),'0')::NUMERIC ),'  ',' ')))::VARCHAR;
        END IF;
    ELSE
        IF split_part(importe::VARCHAR, '.',1) = '1' THEN
            RETURN concat_ws(' ','Pesos uno con', lower(replace(fu_numero_letras(coalesce(split_part(importe::VARCHAR, '.',2),'0')::NUMERIC),'  ',' ')) , 'centavos')::VARCHAR;
        ELSE
            RETURN concat_ws(' ','Pesos', lower(replace(fu_numero_letras(coalesce(split_part(importe::VARCHAR, '.',1),'0')::NUMERIC ),'  ',' ')), 'con' , lower(replace(fu_numero_letras(coalesce(split_part(importe::VARCHAR, '.',2),'0')::NUMERIC),'  ',' ')) , 'centavos')::VARCHAR;
        END IF;
    END IF;

END;
$BODY$
    LANGUAGE plpgsql VOLATILE
                     COST 100;
ALTER FUNCTION getImporteEnLetras(numeric)
    OWNER TO postgres;
