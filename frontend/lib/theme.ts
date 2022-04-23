import { extendTheme } from '@chakra-ui/react';
import { mode } from '@chakra-ui/theme-tools';

const styles = {
    global: (props: any) => ({
        "body": {
          backgroundColor: mode("gray.50", "gray.900")(props),
        },
    }),
};

const components = {
    // ButtonGroup: {
    //     defaultProps: {
    //         colorScheme: "twitter"
    //     },
    // },
    // Button: {
    //     defaultProps: {
    //         colorScheme: "twitter"
    //     },
    // },
    // IconButton: {
    //     defaultProps: {
    //         colorScheme: "twitter"
    //     },
    // },
};

const theme = extendTheme({
    styles,
    components,
});

export default theme;