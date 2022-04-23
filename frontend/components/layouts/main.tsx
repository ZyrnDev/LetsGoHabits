import React, { FC } from 'react';
import { Box, BoxProps, Center, CenterProps, Text, useColorModeValue } from '@chakra-ui/react';
import { Router } from 'next/router';
import Navbar from '@components/navbar';

export const padding = "0.5rem";
export const heading_height = "3.5rem";
export const footer_height = "2rem";
export const content_min_height = `100vh - ${heading_height} - ${footer_height} -  4 * ${padding}`;

const Layout: FC<{ router: Router }> = ({ router, children }) => (
  <>
    <Header p={padding} h={heading_height} />
    <Box as="main" minH={`calc(${content_min_height})`} p={padding}>
      {children}
    </Box>
  </>
);

export const Header: FC<BoxProps> = ({ ...props }) => (
  <Box as="header" top="0" w="100%" position="sticky" zIndex="1" backdropFilter="blur(10px)" {...props}>
    <Navbar w="100%" />
  </Box>
);

export const Footer: FC<CenterProps> = ({ ...props }) => (
  <Center as="footer" w="100%" {...props}>
    made with <Text as="span" color="dark.red" p="0 0.2rem">❤</Text> by zyrn
  </Center>
);


export default Layout;